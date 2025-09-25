package progress

import "regexp"

func GetNumberLineGoto(line string) string {
	rgx, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return ""
	}
	return rgx.ReplaceAllString(line, "")
}

func Percent(currentNumberLine, totalLines int) float64 {
	return float64(currentNumberLine*100.0) / float64(totalLines)
}

func LinesToChangePercentagePoint(currentLine, totalLines int) int {
	start := currentLine
	linesToChangePercentage := -1
	percentageWithCurrentLine := int(Percent(currentLine, totalLines))
	for {
		currentLine++
		nextPercentage := int(Percent(currentLine, totalLines))
		if nextPercentage > percentageWithCurrentLine {
			linesToChangePercentage = currentLine
			break
		}
	}

	return linesToChangePercentage - start
}

func GetPercentage(currentPosition int, fileContent *[]string) float64 {
	percent := float64(currentPosition) * 100.00
	percent = percent / float64(len(*fileContent))
	return percent
}
