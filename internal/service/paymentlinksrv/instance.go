package paymentlinksrv

import "time"

var NowFunc = time.Now

type Impl struct {
	Now func() time.Time
}

func New() PaymentLinkService {
	return &Impl{
		Now: NowFunc,
	}
}
