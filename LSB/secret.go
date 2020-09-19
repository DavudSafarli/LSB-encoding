package LSB

import "strconv"

type Secret struct {
	Text string
}

// getMetadata returns metadata, which is the length of text as binary string
func (s Secret) GetMetadata() string {
	return strconv.FormatInt(int64(len(s.Text)), 2)
}

func (s Secret) GetMetadataAsBitSlice(desiredLength int) []bool {
	metadata := s.GetMetadata()

	metadataBitSlice := make([]bool, desiredLength-len(metadata), desiredLength)
	for _, char := range s.GetMetadata() {
		if char == '1' {
			metadataBitSlice = append(metadataBitSlice, true)
			continue
		}
		metadataBitSlice = append(metadataBitSlice, false)
	}
	return metadataBitSlice
}

func (s Secret) TextAsASCIIBits() []bool {
	textAsAsciiBits := make([]bool, 0, len(s.Text)*8)
	for _, asciiCode := range s.Text {
		asciiCodeInBinary := strconv.FormatInt(int64(asciiCode), 2)
		asciiCodeAsBitSlice := make([]bool, 8-len(asciiCodeInBinary), 8)

		for _, char := range asciiCodeInBinary {
			if char == '1' {
				asciiCodeAsBitSlice = append(asciiCodeAsBitSlice, true)
				continue
			}
			asciiCodeAsBitSlice = append(asciiCodeAsBitSlice, false)
		}
		textAsAsciiBits = append(textAsAsciiBits, asciiCodeAsBitSlice...)
	}
	return textAsAsciiBits
}

func ConvertAsciiBitsToText(arr []bool) string {
	text := ""
	if len(arr)%8 != 0 {
		panic("something went wrong while converting text to ascii array. each character should have been 8 bit long")
	}

	asciiInBinary := ""
	for i, b := range arr {
		if b {
			asciiInBinary += "1"
		} else {
			asciiInBinary += "0"
		}
		if (i+1)%8 == 0 {
			asciiCode, _ := strconv.ParseInt(asciiInBinary, 2, 64)
			text += string(rune(asciiCode))
			asciiInBinary = ""
		}
	}

	return text
}
