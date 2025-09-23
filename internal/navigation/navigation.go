package navigation

import "textreader/internal/model"

func UpdateRangesUp() {
	if model.From <= 0 {
		return
	}

	if model.From > 0 {
		model.From--
	}

	model.To--
}

func UpdateRangesReferenceUp() {
	if model.FromForReferences <= 0 {
		return
	}

	if model.FromForReferences > 0 {
		model.FromForReferences--
	}

	model.ToReferences--
}

func UpdateRangesDown() {
	if model.From < len(model.FileContent) {
		model.From++
	}

	if model.To >= len(model.FileContent) {
		return
	}

	if model.To < len(model.FileContent) {
		model.To++
	}
}

func UpdateRangesReferenceDown() {
	if model.FromForReferences < len(model.References) {
		model.FromForReferences++
	}

	if model.ToReferences >= len(model.References) {
		return
	}

	if model.ToReferences < len(model.References) {
		model.ToReferences++
	}
}
