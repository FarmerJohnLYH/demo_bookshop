// Code generated by Kitex v0.4.3. DO NOT EDIT.

package itemservice

import (
	"context"
	item "github.com/cloudwego/biz-demo/book-shop/kitex_gen/cwg/bookshop/item"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	Add(ctx context.Context, req *item.AddReq, callOptions ...callopt.Option) (r *item.AddResp, err error)
	Edit(ctx context.Context, req *item.EditReq, callOptions ...callopt.Option) (r *item.EditResp, err error)
	Delete(ctx context.Context, req *item.DeleteReq, callOptions ...callopt.Option) (r *item.DeleteResp, err error)
	Online(ctx context.Context, req *item.OnlineReq, callOptions ...callopt.Option) (r *item.OnlineResp, err error)
	Offline(ctx context.Context, req *item.OfflineReq, callOptions ...callopt.Option) (r *item.OfflineResp, err error)
	Get(ctx context.Context, req *item.GetReq, callOptions ...callopt.Option) (r *item.GetResp, err error)
	MGet2C(ctx context.Context, req *item.MGet2CReq, callOptions ...callopt.Option) (r *item.MGet2CResp, err error)
	Search(ctx context.Context, req *item.SearchReq, callOptions ...callopt.Option) (r *item.SearchResp, err error)
	List(ctx context.Context, req *item.ListReq, callOptions ...callopt.Option) (r *item.ListResp, err error)
	DecrStock(ctx context.Context, req *item.DecrStockReq, callOptions ...callopt.Option) (r *item.DecrStockResp, err error)
	DecrStockRevert(ctx context.Context, req *item.DecrStockReq, callOptions ...callopt.Option) (r *item.DecrStockResp, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kItemServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kItemServiceClient struct {
	*kClient
}

func (p *kItemServiceClient) Add(ctx context.Context, req *item.AddReq, callOptions ...callopt.Option) (r *item.AddResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Add(ctx, req)
}

func (p *kItemServiceClient) Edit(ctx context.Context, req *item.EditReq, callOptions ...callopt.Option) (r *item.EditResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Edit(ctx, req)
}

func (p *kItemServiceClient) Delete(ctx context.Context, req *item.DeleteReq, callOptions ...callopt.Option) (r *item.DeleteResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Delete(ctx, req)
}

func (p *kItemServiceClient) Online(ctx context.Context, req *item.OnlineReq, callOptions ...callopt.Option) (r *item.OnlineResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Online(ctx, req)
}

func (p *kItemServiceClient) Offline(ctx context.Context, req *item.OfflineReq, callOptions ...callopt.Option) (r *item.OfflineResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Offline(ctx, req)
}

func (p *kItemServiceClient) Get(ctx context.Context, req *item.GetReq, callOptions ...callopt.Option) (r *item.GetResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Get(ctx, req)
}

func (p *kItemServiceClient) MGet2C(ctx context.Context, req *item.MGet2CReq, callOptions ...callopt.Option) (r *item.MGet2CResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.MGet2C(ctx, req)
}

func (p *kItemServiceClient) Search(ctx context.Context, req *item.SearchReq, callOptions ...callopt.Option) (r *item.SearchResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Search(ctx, req)
}

func (p *kItemServiceClient) List(ctx context.Context, req *item.ListReq, callOptions ...callopt.Option) (r *item.ListResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.List(ctx, req)
}

func (p *kItemServiceClient) DecrStock(ctx context.Context, req *item.DecrStockReq, callOptions ...callopt.Option) (r *item.DecrStockResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.DecrStock(ctx, req)
}

func (p *kItemServiceClient) DecrStockRevert(ctx context.Context, req *item.DecrStockReq, callOptions ...callopt.Option) (r *item.DecrStockResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.DecrStockRevert(ctx, req)
}
