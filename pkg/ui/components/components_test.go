package components

import (
	"strings"
	"testing"
)

func TestNewButton(t *testing.T) {
	btn := NewButton("Click me")

	if btn.Content != "Click me" {
		t.Errorf("Expected content 'Click me', got '%s'", btn.Content)
	}

	if btn.Variant != ButtonPrimary {
		t.Errorf("Expected primary variant, got %v", btn.Variant)
	}

	if btn.Size != ButtonSizeMedium {
		t.Errorf("Expected medium size, got %v", btn.Size)
	}

	if btn.Disabled {
		t.Error("Expected button to not be disabled")
	}
}

func TestButtonRender(t *testing.T) {
	btn := NewButton("Test")
	rendered := btn.Render()

	if rendered == "" {
		t.Error("Button render returned empty string")
	}

	if !strings.Contains(rendered, "Test") {
		t.Error("Button render should contain the button text")
	}
}

func TestButtonWithIcon(t *testing.T) {
	btn := NewButton("Save")
	btn.Icon = "💾"
	rendered := btn.Render()

	if !strings.Contains(rendered, "💾") {
		t.Error("Button render should contain the icon")
	}

	if !strings.Contains(rendered, "Save") {
		t.Error("Button render should contain the text")
	}
}

func TestButtonVariants(t *testing.T) {
	variants := []ButtonVariant{
		ButtonPrimary,
		ButtonSecondary,
		ButtonDanger,
		ButtonGhost,
		ButtonLink,
	}

	for _, variant := range variants {
		btn := NewButton("Test")
		btn.Variant = variant
		rendered := btn.Render()

		if rendered == "" {
			t.Errorf("Button variant %s rendered empty", variant)
		}
	}
}

func TestButtonSizes(t *testing.T) {
	sizes := []ButtonSize{
		ButtonSizeSmall,
		ButtonSizeMedium,
		ButtonSizeLarge,
	}

	for _, size := range sizes {
		btn := NewButton("Test")
		btn.Size = size
		rendered := btn.Render()

		if rendered == "" {
			t.Errorf("Button size %s rendered empty", size)
		}
	}
}

func TestButtonDisabled(t *testing.T) {
	btn := NewButton("Disabled")
	btn.Disabled = true
	rendered := btn.Render()

	if rendered == "" {
		t.Error("Disabled button should still render")
	}
}

func TestNewCard(t *testing.T) {
	card := NewCard("Title", "Content")

	if card.Title != "Title" {
		t.Errorf("Expected title 'Title', got '%s'", card.Title)
	}

	if card.Content != "Content" {
		t.Errorf("Expected content 'Content', got '%s'", card.Content)
	}

	if !card.Elevated {
		t.Error("Expected card to be elevated by default")
	}

	if !card.Bordered {
		t.Error("Expected card to be bordered by default")
	}
}

func TestCardRender(t *testing.T) {
	card := NewCard("Test Title", "Test Content")
	rendered := card.Render()

	if rendered == "" {
		t.Error("Card render returned empty string")
	}

	if !strings.Contains(rendered, "Test Title") {
		t.Error("Card render should contain the title")
	}

	if !strings.Contains(rendered, "Test Content") {
		t.Error("Card render should contain the content")
	}
}

func TestCardWithFooter(t *testing.T) {
	card := NewCard("Title", "Content")
	card.Footer = "Footer text"
	rendered := card.Render()

	if !strings.Contains(rendered, "Footer text") {
		t.Error("Card render should contain the footer")
	}
}

func TestNewBadge(t *testing.T) {
	badge := NewBadge("NEW")

	if badge.Label != "NEW" {
		t.Errorf("Expected label 'NEW', got '%s'", badge.Label)
	}

	if badge.Variant != BadgeDefault {
		t.Errorf("Expected default variant, got %v", badge.Variant)
	}
}

func TestBadgeRender(t *testing.T) {
	badge := NewBadge("Test")
	rendered := badge.Render()

	if rendered == "" {
		t.Error("Badge render returned empty string")
	}

	if !strings.Contains(rendered, "Test") {
		t.Error("Badge render should contain the label")
	}
}

func TestBadgeVariants(t *testing.T) {
	variants := []BadgeVariant{
		BadgeDefault,
		BadgeSuccess,
		BadgeWarning,
		BadgeError,
		BadgeInfo,
	}

	for _, variant := range variants {
		badge := NewBadge("Test")
		badge.Variant = variant
		rendered := badge.Render()

		if rendered == "" {
			t.Errorf("Badge variant %s rendered empty", variant)
		}
	}
}

func TestNewStack(t *testing.T) {
	stack := NewStack(StackVertical)

	if stack.Direction != StackVertical {
		t.Errorf("Expected vertical direction, got %v", stack.Direction)
	}

	if len(stack.Items) != 0 {
		t.Error("New stack should have no items")
	}
}

func TestStackAdd(t *testing.T) {
	stack := NewStack(StackVertical)
	stack.Add("Item 1").Add("Item 2")

	if len(stack.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(stack.Items))
	}

	if stack.Items[0] != "Item 1" {
		t.Errorf("Expected first item 'Item 1', got '%s'", stack.Items[0])
	}
}

func TestStackRenderVertical(t *testing.T) {
	stack := NewStack(StackVertical)
	stack.Add("Line 1").Add("Line 2").Add("Line 3")
	rendered := stack.Render()

	if rendered == "" {
		t.Error("Stack render returned empty string")
	}

	if !strings.Contains(rendered, "Line 1") {
		t.Error("Stack should contain Line 1")
	}
	if !strings.Contains(rendered, "Line 2") {
		t.Error("Stack should contain Line 2")
	}
	if !strings.Contains(rendered, "Line 3") {
		t.Error("Stack should contain Line 3")
	}
}

func TestStackRenderHorizontal(t *testing.T) {
	stack := NewStack(StackHorizontal)
	stack.Add("Item1").Add("Item2")
	rendered := stack.Render()

	if rendered == "" {
		t.Error("Stack render returned empty string")
	}

	if !strings.Contains(rendered, "Item1") || !strings.Contains(rendered, "Item2") {
		t.Error("Stack should contain all items")
	}
}

func TestStackEmpty(t *testing.T) {
	stack := NewStack(StackVertical)
	rendered := stack.Render()

	if rendered != "" {
		t.Error("Empty stack should render empty string")
	}
}

func TestNewDivider(t *testing.T) {
	divider := NewDivider(10)

	if divider.Length != 10 {
		t.Errorf("Expected length 10, got %d", divider.Length)
	}

	if divider.Char != "─" {
		t.Errorf("Expected char '─', got '%s'", divider.Char)
	}
}

func TestDividerRender(t *testing.T) {
	divider := NewDivider(5)
	rendered := divider.Render()

	if rendered == "" {
		t.Error("Divider render returned empty string")
	}

	// The rendered output contains ANSI codes, so just check it's not empty
	if len(rendered) < 5 {
		t.Error("Divider render should have content")
	}
}

func TestNewStatusIndicator(t *testing.T) {
	status := NewStatusIndicator("Running", StatusTypeInfo)

	if status.Label != "Running" {
		t.Errorf("Expected label 'Running', got '%s'", status.Label)
	}

	if status.Status != StatusTypeInfo {
		t.Errorf("Expected info status, got %v", status.Status)
	}

	if !status.ShowIcon {
		t.Error("Expected ShowIcon to be true by default")
	}
}

func TestStatusIndicatorRender(t *testing.T) {
	status := NewStatusIndicator("Test", StatusTypeSuccess)
	rendered := status.Render()

	if rendered == "" {
		t.Error("Status indicator render returned empty string")
	}

	if !strings.Contains(rendered, "Test") {
		t.Error("Status indicator should contain the label")
	}
}

func TestStatusIndicatorTypes(t *testing.T) {
	types := []StatusType{
		StatusTypeSuccess,
		StatusTypeWarning,
		StatusTypeError,
		StatusTypeInfo,
		StatusTypeNeutral,
	}

	for _, statusType := range types {
		status := NewStatusIndicator("Test", statusType)
		rendered := status.Render()

		if rendered == "" {
			t.Errorf("Status type %s rendered empty", statusType)
		}
	}
}

func TestStatusIndicatorWithoutIcon(t *testing.T) {
	status := NewStatusIndicator("No Icon", StatusTypeInfo)
	status.ShowIcon = false
	rendered := status.Render()

	if !strings.Contains(rendered, "No Icon") {
		t.Error("Status indicator should contain label even without icon")
	}
}

func TestComponentComposition(t *testing.T) {
	// Test that components can be composed together
	badge := NewBadge("NEW")
	button := NewButton("Click me")
	button.Icon = badge.Render() // Compose badge into button

	rendered := button.Render()
	if rendered == "" {
		t.Error("Composed component should render")
	}
}
