// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"deutsch/internal/code"
	"deutsch/internal/config"
	"deutsch/internal/handler"
	"deutsch/internal/svc"
	"deutsch/internal/types"
	"deutsch/model/gormdb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/deutsch.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if c.MySQL.DataSource == "" {
		log.Fatal("MySQL DataSource is required")
	}
	if err := gormdb.InitDB(c.MySQL.DataSource); err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, any) {
		var ce *code.CodeError
		if errors.As(err, &ce) {
			return http.StatusOK, &types.Base{Code: int(ce.Code), Msg: ce.Error()}
		}
		return http.StatusInternalServerError, &types.Base{
			Code: int(code.CodeInternalServerError),
			Msg:  err.Error(),
		}
	})

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx, err := svc.NewServiceContext(c)
	if err != nil {
		log.Fatalf("failed to initialize service context: %v", err)
	}
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
