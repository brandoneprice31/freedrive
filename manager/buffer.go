package manager

import (
	"encoding/json"
	"sort"
	"sync"

	"github.com/brandoneprice31/freedrive/service"
)

type (
	Buffer struct {
		service service.Service

		serviceData []service.ServiceData
		wg          *sync.WaitGroup
		sdMutex     *sync.Mutex
		serviceErrs []error
		seMutex     *sync.Mutex

		concCount      int
		concCountMutex *sync.Mutex
		concCountCond  *sync.Cond

		buf     []obj
		bufSize int

		dataTag int

		downloaded map[int][]obj
		objMutex   *sync.Mutex
	}

	obj struct {
		Data []byte `json:"data"`
		Tag  int    `json:"tag"`
		Iter int    `json:"iter"`
	}

	objs []obj
)

func (b *Buffer) ServiceTye() service.ServiceType {
	return b.service.Type()
}

func NewBackupBuffer(s service.Service) (*Buffer, error) {
	err := s.NewBackup()
	if err != nil {
		return nil, err
	}

	concCountMutex := &sync.Mutex{}

	return &Buffer{
		service:        s,
		wg:             &sync.WaitGroup{},
		sdMutex:        &sync.Mutex{},
		seMutex:        &sync.Mutex{},
		concCountMutex: concCountMutex,
		concCountCond:  sync.NewCond(concCountMutex),
		dataTag:        0,
	}, nil
}

func (b *Buffer) Save(data []byte) error {
	err := b.save(obj{
		Data: data,
		Tag:  b.dataTag,
		Iter: 0,
	})
	if err != nil {
		return err
	}

	b.dataTag++
	return nil
}

func (b *Buffer) save(o obj) error {
	if len(o.Data) == 0 {
		return nil
	}

	if b.bufSize == b.service.BufferSize() {
		err := b.uploadBuf()
		if err != nil {
			return err
		}

		b.buf = []obj{}
		b.bufSize = 0
	}

	spaceLeft := b.service.BufferSize() - b.bufSize

	if spaceLeft >= len(o.Data) {
		b.buf = append(b.buf, o)
		b.bufSize += len(o.Data)
		return nil
	}

	err := b.save(obj{
		Data: o.Data[:spaceLeft],
		Tag:  o.Tag,
		Iter: o.Iter,
	})
	if err != nil {
		return err
	}

	return b.save(obj{
		Data: o.Data[spaceLeft:],
		Tag:  o.Tag,
		Iter: o.Iter + 1,
	})
}

func (b *Buffer) uploadBuf() error {
	data, err := json.Marshal(b.buf)
	if err != nil {
		return err
	}

	b.wg.Add(1)
	go func() {
		b.concCountMutex.Lock()
		for {
			if b.concCount >= b.service.MaxThreads() {
				b.concCountCond.Wait()
			} else {
				break
			}
		}
		b.concCount++
		b.concCountMutex.Unlock()

		sd, err := b.service.Upload(data)
		if err != nil {
			b.appendServiceErr(err)
		} else {
			b.appendServiceData(sd)
		}

		b.concCountMutex.Lock()
		b.concCount--
		b.concCountCond.Signal()
		b.concCountMutex.Unlock()

		b.wg.Done()
	}()

	return nil
}

func (b *Buffer) appendServiceData(sd *service.ServiceData) {
	b.sdMutex.Lock()
	b.serviceData = append(b.serviceData, *sd)
	b.sdMutex.Unlock()
}

func (b *Buffer) appendServiceErr(err error) {
	b.seMutex.Lock()
	b.serviceErrs = append(b.serviceErrs, err)
	b.seMutex.Unlock()
}

func (b *Buffer) FlushBackup() ([]service.ServiceData, error) {
	if b.bufSize > 0 {
		err := b.uploadBuf()
		if err != nil {
			return nil, err
		}
	}

	b.wg.Wait()

	if len(b.serviceErrs) > 0 {
		return nil, b.serviceErrs[0]
	}

	return b.serviceData, nil
}

func NewDownloadBuffer(s service.Service) (*Buffer, error) {
	err := s.NewDownload()
	if err != nil {
		return nil, err
	}

	concCountMutex := &sync.Mutex{}

	return &Buffer{
		service:        s,
		downloaded:     make(map[int][]obj),
		objMutex:       &sync.Mutex{},
		seMutex:        &sync.Mutex{},
		concCountMutex: concCountMutex,
		concCountCond:  sync.NewCond(concCountMutex),
		wg:             &sync.WaitGroup{},
	}, nil
}

func (b *Buffer) Download(sd service.ServiceData) error {
	b.wg.Add(1)
	go func() {
		b.concCountMutex.Lock()
		for {
			if b.concCount >= b.service.MaxThreads() {
				b.concCountCond.Wait()
			} else {
				break
			}
		}
		b.concCount++
		b.concCountMutex.Unlock()

		data, err := b.service.Download(sd)
		if err != nil {
			b.appendServiceErr(err)

		} else {
			var oo []obj
			err = json.Unmarshal(data, &oo)
			if err != nil {
				b.appendServiceErr(err)

			} else {
				for _, o := range oo {
					b.appendDownloaded(o)
				}
			}
		}

		b.concCountMutex.Lock()
		b.concCount--
		b.concCountCond.Signal()
		b.concCountMutex.Unlock()

		b.wg.Done()
	}()

	return nil
}

func (b *Buffer) appendDownloaded(o obj) {
	b.objMutex.Lock()
	b.downloaded[o.Tag] = append(b.downloaded[o.Tag], o)
	b.objMutex.Unlock()
}

func (b *Buffer) FlushDownload() ([][]byte, error) {
	b.wg.Wait()

	if len(b.serviceErrs) > 0 {
		return nil, b.serviceErrs[0]
	}

	dd := [][]byte{}
	for _, oo := range b.downloaded {
		sort.Sort(objs(oo))

		d := []byte{}
		for _, o := range oo {
			d = append(d, o.Data...)
		}

		dd = append(dd, d)
	}

	return dd, nil
}

func (oo objs) Len() int {
	return len(oo)
}

func (oo objs) Less(i, j int) bool {
	return oo[i].Iter < oo[j].Iter
}

func (oo objs) Swap(i, j int) {
	oo[i], oo[j] = oo[j], oo[i]
}
