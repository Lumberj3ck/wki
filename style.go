package main

import "github.com/charmbracelet/lipgloss"

var (
	ACCENT_COLOR = "#04B575"
	titleStyle   = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()

	linkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ACCENT_COLOR)).
			Render

	listArticleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ACCENT_COLOR)).
				Render
	articleDescriptionStyle = lipgloss.NewStyle().
				Bold(true).
				Underline(true).
				Render
	articleBoldedStyle = lipgloss.NewStyle().
				Bold(true).
				Render
	articleItalicStyle = lipgloss.NewStyle().
				Italic(true).
				Render
	articleBoldedItalicStyle = lipgloss.NewStyle().
					Bold(true).
					Italic(true).
					Render
	noteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080")).
			Render
)
