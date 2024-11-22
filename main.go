package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"codeberg.org/withlove/mono/pgdb"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
)

func main() {
	ctx := context.Background()
	cfg := lo.Must(pgdb.ConfigFromEnv())
	pool := lo.Must(pgdb.OpenPool(cfg, zerolog.New(os.Stderr).Level(zerolog.WarnLevel).With().Timestamp().Logger()))
	db := pgdb.ToSQLX(pool)
	fmt.Printf("pool opened\n\n")

	ids := make([]int64, 0)
	lo.Must0(db.SelectContext(
		ctx,
		&ids,
		`select id from public.temp order by id desc`,
	))
	fmt.Println(ids)

	time.Sleep(30 * time.Second)

}
