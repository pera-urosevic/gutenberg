package stamps

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"somnusalis.org/gutenberg/system"

	"gorm.io/gorm"
)

type Stamp struct {
	Id         uint
	Image      []byte
	Name       string
	Country    string
	Colnect    string
	Issued     string
	Duplicates int
}

func Stamps() {
	where := "false"
	if len(os.Args) > 2 {
		where = strings.Join(os.Args[2:], " ")
	}

	system.Log("Stamps", STAMPS_DB)
	db, err := gorm.Open(sqlite.Open(STAMPS_DB), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to stamps database")
	}

	system.Log("Where", where)
	var stamps []Stamp
	db.Where(where).Order("country, issued").Find(&stamps)
	system.Log("Found", len(stamps))

	m := getMaroto(stamps)
	document, err := m.Generate()
	if err != nil {
		log.Fatal(err.Error())
	}

	system.Log("Save", STAMPS_PDF)
	err = document.Save(STAMPS_PDF)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getMaroto(stamps []Stamp) core.Maroto {
	cfg := config.NewBuilder().
		WithPageNumber("Stamps - Duplicates, page {current} of {total}", props.LeftBottom).
		WithMargins(5, 5, 5).
		WithOrientation(orientation.Vertical).
		// WithDebug(true).
		Build()

	mrt := maroto.New(cfg)
	m := maroto.NewMetricsDecorator(mrt)

	m.AddRow(3)
	for i, stamp := range stamps {
		url := "https://colnect.com/en/stamps/stamp/" + stamp.Colnect
		m.AddRow(24,
			image.NewFromBytesCol(2, stamp.Image, extension.Jpg, props.Rect{Center: true}),
			col.New(10).Add(
				text.New(
					strconv.Itoa(i+1)+". "+stamp.Name+" ("+strconv.Itoa(stamp.Duplicates)+")",
					props.Text{
						Top:       0,
						Left:      2,
						Size:      12,
						Hyperlink: &url,
					},
				),
				text.New(
					stamp.Country,
					props.Text{
						Top:  6,
						Left: 2,
					},
				),
				text.New(
					stamp.Issued,
					props.Text{
						Top:  11,
						Left: 2,
					},
				),
			),
		)
		m.AddRow(3)
	}

	return m
}
