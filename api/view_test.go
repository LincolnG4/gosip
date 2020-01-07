package api

import (
	"testing"
)

func TestView(t *testing.T) {
	checkClient(t)

	web := NewSP(spClient).Web()
	listURI := getRelativeURL(spClient.AuthCnfg.GetSiteURL()) + "/Shared%20Documents"

	t.Run("Conf", func(t *testing.T) {
		f := web.Fields().GetByID("")
		hs := map[string]*RequestConfig{
			"nometadata":      HeadersPresets.Nometadata,
			"minimalmetadata": HeadersPresets.Minimalmetadata,
			"verbose":         HeadersPresets.Verbose,
		}
		for key, preset := range hs {
			g := f.Conf(preset)
			if g.config != preset {
				t.Errorf("can't %v config", key)
			}
		}
	})

	t.Run("Get", func(t *testing.T) {
		data, err := web.GetList(listURI).Views().DefaultView().Get()
		if err != nil {
			t.Error(err)
		}
		if data.Data().ID == "" {
			t.Error("can't unmarshal data")
		}
	})

	t.Run("Modifiers", func(t *testing.T) {
		mods := web.GetList(listURI).Views().DefaultView().
			Select("*").Expand("*").modifiers
		if mods == nil || len(mods.mods) != 2 {
			t.Error("can't add modifiers")
		}
	})

	t.Run("FromURL", func(t *testing.T) {
		viewR, err := web.GetList(listURI).Views().DefaultView().Get()
		if err != nil {
			t.Error(err)
		}
		entityURL := ExtractEntityURI(viewR)
		if entityURL == "" {
			t.Error("can't extract entity URL")
		}
		view1, err := web.GetList(listURI).Views().GetByID("").FromURL(entityURL).Get()
		if err != nil {
			t.Error(err)
		}
		if viewR.Data().ID != view1.Data().ID {
			t.Error("can't get view from entity URL")
		}
	})

	// ToDo:
	// Update
	// Delete
	// Recycle

}
