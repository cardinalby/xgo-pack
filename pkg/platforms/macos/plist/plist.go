package plist

import (
	"encoding/xml"
)

type plist struct {
	XMLName    xml.Name   `xml:"plist"`
	Version    string     `xml:"version,attr"`
	Dictionary dictionary `xml:"dict"`
}

type dictionary struct {
	Entries []dictEntry
}

type dictEntry struct {
	Key   string `xml:"key"`
	Value string `xml:"string"`
}

func (d dictionary) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for _, entry := range d.Entries {
		if err := e.EncodeElement(entry.Key, xml.StartElement{Name: xml.Name{Local: "key"}}); err != nil {
			return err
		}
		if err := e.EncodeElement(entry.Value, xml.StartElement{Name: xml.Name{Local: "string"}}); err != nil {
			return err
		}
	}

	if err := e.EncodeToken(start.End()); err != nil {
		return err
	}
	return nil
}

//goland:noinspection HttpUrlsUsage
const plistHeader = `<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">`

func GetPlist(config Properties) ([]byte, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	plist := plist{
		Version: "1.0",
	}

	if config.ProductName != "" {
		plist.Dictionary.Entries = append(plist.Dictionary.Entries, dictEntry{
			Key:   "CFBundleName",
			Value: config.ProductName,
		})
	}

	if config.BinName != "" {
		plist.Dictionary.Entries = append(plist.Dictionary.Entries, dictEntry{
			Key:   "CFBundleExecutable",
			Value: config.BinName,
		})
	}

	if config.Identifier != "" {
		plist.Dictionary.Entries = append(plist.Dictionary.Entries, dictEntry{
			Key:   "CFBundleIdentifier",
			Value: config.Identifier,
		})
	}

	if config.ProductVersion != "" {
		plist.Dictionary.Entries = append(plist.Dictionary.Entries, dictEntry{
			Key:   "CFBundleVersion",
			Value: config.ProductVersion,
		})
		plist.Dictionary.Entries = append(plist.Dictionary.Entries, dictEntry{
			Key:   "CFBundleShortVersionString",
			Value: config.ProductVersion,
		})
	}

	if config.IconFile != "" {
		plist.Dictionary.Entries = append(plist.Dictionary.Entries, dictEntry{
			Key:   "CFBundleIconFile",
			Value: config.IconFile,
		})
	}

	if config.IsHighResolutionCapable {
		plist.Dictionary.Entries = append(plist.Dictionary.Entries, dictEntry{
			Key:   "NSHighResolutionCapable",
			Value: "true",
		})
	}

	if config.IsHideInDock {
		plist.Dictionary.Entries = append(plist.Dictionary.Entries, dictEntry{
			Key:   "LSUIElement",
			Value: "1",
		})
	}

	if config.CopyRight != "" {
		plist.Dictionary.Entries = append(plist.Dictionary.Entries, dictEntry{
			Key:   "CFBundleGetInfoString",
			Value: config.CopyRight,
		})
	}

	data, err := xml.MarshalIndent(plist, "", "    ")
	if err != nil {
		return nil, err
	}
	return []byte(plistHeader + "\n" + string(data)), nil
}
