package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestItems(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	newListTitle := strings.Replace(uuid.New().String(), "-", "", -1)
	if _, err := web.Lists().Add(newListTitle, nil); err != nil {
		t.Error(err)
	}
	list := web.Lists().GetByTitle(newListTitle)
	entType, err := list.GetEntityType()
	if err != nil {
		t.Error(err)
	}

	t.Run("AddWithoutMetadataType", func(t *testing.T) {
		body := []byte(`{"Title":"Item"}`)
		if _, err := list.Items().Add(body); err != nil {
			t.Error(err)
		}
	})

	t.Run("AddResponse", func(t *testing.T) {
		body := []byte(`{"Title":"Item"}`)
		item, err := list.Items().Add(body)
		if err != nil {
			t.Error(err)
		}
		if item.Data().ID == 0 {
			t.Error("can't get item properly")
		}
	})

	t.Run("AddSeries", func(t *testing.T) {
		for i := 1; i < 10; i++ {
			metadata := make(map[string]interface{})
			metadata["__metadata"] = map[string]string{"type": entType}
			metadata["Title"] = fmt.Sprintf("Item %d", i)
			body, _ := json.Marshal(metadata)
			if _, err := list.Items().Add(body); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("AddValidate", func(t *testing.T) {
		if envCode == "2013" {
			t.Skip("is not supported with SP 2013")
		}

		options := &ValidateAddOptions{NewDocumentUpdate: true, CheckInComment: "test"}
		data := map[string]string{"Title": "New item"}
		if _, err := list.Items().AddValidate(data, options); err != nil {
			t.Error(err)
		}
	})

	t.Run("AddValidateWithPath", func(t *testing.T) {
		if envCode == "2013" {
			t.Skip("is not supported with SP 2013")
		}
		// doesn't work anymore in SPO, an item can't be created in a folder which is not item-folder
		// if _, err := list.RootFolder().Folders().Add("subfolder"); err != nil {
		// 	t.Error(err)
		// }

		folderName := "subfolder"

		if _, err := list.Update([]byte(`{ "EnableFolderCreation": true }`)); err != nil {
			t.Error(err)
		}
		ff, err := list.Items().AddValidate(map[string]string{
			"Title":         folderName,
			"FileLeafRef":   folderName,
			"ContentType":   "Folder",
			"ContentTypeId": "0x0120",
		}, nil)
		if err != nil {
			t.Error(err)
		}
		if _, err := list.Items().GetByID(ff.ID()).Update([]byte(`{ "FileLeafRef": "` + folderName + `" }`)); err != nil {
			t.Error(err)
		}

		options := &ValidateAddOptions{NewDocumentUpdate: true, CheckInComment: "test"}
		options.DecodedPath = "Lists/" + newListTitle + "/" + folderName
		data := map[string]string{"Title": "New item in folder"}
		if _, err := list.Items().AddValidate(data, options); err != nil {
			t.Error(err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		items, err := list.Items().Top(100).OrderBy("Title", false).Get()
		if err != nil {
			t.Error(err)
		}
		if len(items.Data()) == 0 {
			t.Error("can't get items properly")
		}
		if items.Data()[0].Data().ID == 0 {
			t.Error("can't get items properly")
		}
		if bytes.Compare(items, items.Normalized()) == -1 {
			t.Error("wrong response normalization")
		}
		if len(items.ToMap()) == 0 {
			t.Error("can't map items properly")
		}
		if items.ToMap()[0]["ID"] == 0 {
			t.Error("can't map items properly")
		}
	})

	t.Run("GetPaged", func(t *testing.T) {
		paged, err := list.Items().Top(5).GetPaged()
		if err != nil {
			t.Error(err)
		}
		if len(paged.Items.Data()) == 0 {
			t.Error("can't get items")
		}
		if !paged.HasNextPage() {
			t.Error("can't get next page")
		} else {
			if _, err := paged.GetNextPage(); err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		item, err := list.Items().GetByID(1).Get()
		if err != nil {
			t.Error(err)
		}
		if item.Data().ID == 0 {
			t.Error("can't get item properly")
		}
	})

	t.Run("Get/Unmarshal", func(t *testing.T) {
		item, err := list.Items().GetByID(1).Get()
		if err != nil {
			t.Error(err)
		}
		if item.Data().ID == 0 {
			t.Error("can't get item ID property properly")
		}
		if item.Data().Title == "" {
			t.Error("can't get item Title property properly")
		}
	})

	t.Run("GetByCAML", func(t *testing.T) {
		caml := `
			<View>
				<Query>
					<Where>
						<Eq>
							<FieldRef Name='ID' />
							<Value Type='Number'>3</Value>
						</Eq>
					</Where>
				</Query>
			</View>
		`
		data, err := list.Items().Select("Id").GetByCAML(caml)
		if err != nil {
			t.Error(err)
		}
		if len(data.Data()) != 1 {
			t.Error("incorrect number of items")
		}
		if data.Data()[0].Data().ID != 3 {
			t.Error("incorrect response")
		}
	})

	t.Run("GetByCAMLAdvanced", func(t *testing.T) {
		viewResp, err := list.Views().DefaultView().Select("ListViewXml").Get()
		if err != nil {
			t.Error(err)
		}
		if _, err := list.Items().GetByCAML(viewResp.Data().ListViewXML); err != nil {
			t.Error(err)
		}
	})

	// if err := list.Delete(); err != nil {
	// 	t.Error(err)
	// }

}
