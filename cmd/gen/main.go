package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type cli struct {
	Host     string `env:"DB_HOST" default:"pg.lab.verysmart.house" help:"host"`
	Username string `env:"DB_USER" default:"watchedsky-social" help:"user"`
	Password string `env:"DB_PASSWORD" help:"db password"`
	DB       string `env:"DB_NAME" default:"watchedsky-social"`
}

func main() {
	var args cli
	kong.Parse(&args)

	g := gen.NewGenerator(gen.Config{
		OutPath:           "../../models",
		OutFile:           "gen_query.go",
		ModelPkgPath:      "models",
		WithUnitTest:      false,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.GenerateMode(0),
	})

	db, _ := gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s TimeZone=UTC", args.Host, args.Username, args.Password, args.DB)))
	g.UseDB(db)
	g.WithFileNameStrategy(func(tableName string) string {
		return fmt.Sprintf("%s_model", tableName)
	})

	dataTypeMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"geometry": func(columnType gorm.ColumnType) (dataType string) {
			ct, _ := columnType.ColumnType()
			if strings.Contains(strings.ToLower(ct), "geometry(geometry,4326)") {
				return "orb.Geometry"
			}

			if strings.Contains(strings.ToLower(ct), "geometry(point,4326)") {
				return "orb.Point"
			}

			return "string"
		},
	}

	g.WithDataTypeMap(dataTypeMap)
	g.WithImportPkgPath("github.com/paulmach/orb")
	g.GenerateModel("zones")

	g.Execute()
}
