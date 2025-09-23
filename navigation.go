package main

import "textreader/internal/model"

func updateRangesUp() {
	if model.From <= 0 {
		return
	}

	if model.From > 0 {
		model.From--
	}

	model.To--
}

func updateRangesReferenceUp() {
	if model.FromForReferences <= 0 {
		return
	}

	if model.FromForReferences > 0 {
		model.FromForReferences--
	}

	model.ToReferences--
}

func updateRangesDown() {
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

func updateRangesReferenceDown() {
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
