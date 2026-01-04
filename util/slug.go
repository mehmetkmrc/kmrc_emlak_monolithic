package util

import (
    "strings"
    "unicode"

    "golang.org/x/text/unicode/norm"
)

func Slugify(text string) string {
    text = strings.ToLower(text)

    // Türkçe karakter dönüşümleri
    replacer := strings.NewReplacer(
        "ç", "c",
        "ğ", "g",
        "ı", "i",
        "ö", "o",
        "ş", "s",
        "ü", "u",
    )
    text = replacer.Replace(text)

    // Unicode normalize
    t := norm.NFD.String(text)
    var b strings.Builder
    for _, r := range t {
        if unicode.Is(unicode.Mn, r) {
            continue
        }
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            b.WriteRune(r)
        } else {
            b.WriteRune('-')
        }
    }

    slug := strings.Trim(b.String(), "-")
    slug = strings.ReplaceAll(slug, "--", "-")

    return slug
}
