package entity

import (
	"testing"
	"time"
)

func TestArticle_TableName(t *testing.T) {
	article := Article{}
	expected := "articles"
	actual := article.TableName()
	
	if actual != expected {
		t.Errorf("Expected table name %s, got %s", expected, actual)
	}
}

func TestArticle_Creation(t *testing.T) {
	article := Article{
		Title:     "Test Article",
		Content:   "Test content",
		Author:    "Test Author",
		Published: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	if article.Title != "Test Article" {
		t.Errorf("Expected title %s, got %s", "Test Article", article.Title)
	}
	
	if article.Content != "Test content" {
		t.Errorf("Expected content %s, got %s", "Test content", article.Content)
	}
	
	if article.Author != "Test Author" {
		t.Errorf("Expected author %s, got %s", "Test Author", article.Author)
	}
	
	if !article.Published {
		t.Error("Expected published to be true")
	}
}