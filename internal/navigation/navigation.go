package navigation

import "textreader/internal/model"

func UpdateRangesUp(state *model.AppState) {
	if state.From <= 0 {
		return
	}

	if state.From > 0 {
		state.From--
	}

	state.To--
}

func UpdateRangesReferenceUp(state *model.AppState) {
	if state.FromForReferences <= 0 {
		return
	}

	if state.FromForReferences > 0 {
		state.FromForReferences--
	}

	state.ToReferences--
}

func UpdateRangesDown(state *model.AppState) {
	if state.From < len(state.FileContent) {
		state.From++
	}

	if state.To >= len(state.FileContent) {
		return
	}

	if state.To < len(state.FileContent) {
		state.To++
	}
}

func UpdateRangesReferenceDown(state *model.AppState) {
	if state.FromForReferences < len(state.References) {
		state.FromForReferences++
	}

	if state.ToReferences >= len(state.References) {
		return
	}

	if state.ToReferences < len(state.References) {
		state.ToReferences++
	}
}

func UpdateRangesVocabularyUp(state *model.AppState) {
	if state.FromForVocabulary <= 0 {
		return
	}
	if state.FromForVocabulary > 0 {
		state.FromForVocabulary--
	}
	state.ToForVocabulary--
}

func UpdateRangesVocabularyDown(state *model.AppState) {
	if state.FromForVocabulary < len(state.Vocabulary) {
		state.FromForVocabulary++
	}
	if state.ToForVocabulary >= len(state.Vocabulary) {
		return
	}
	if state.ToForVocabulary < len(state.Vocabulary) {
		state.ToForVocabulary++
	}
}
