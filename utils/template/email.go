package template

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
)

type EmailTemplateData struct {
	Name       string
	VerifyLink string
	ResetLink  string
	Email      string
	AppName    string
}

func getProjectRoot() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(filepath.Dir(filename)))
}

func RenderEmailTemplate(templateName string, data interface{}) (string, error) {
	projectRoot := getProjectRoot()
	templatePath := filepath.Join(projectRoot, "templates", "email", templateName+".html")

	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("templates file not found: %s", templatePath)
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buf.String(), nil
}

// Helper function untuk templates verifikasi email
func RenderVerificationEmail(name, verifyLink string) (string, error) {
	data := EmailTemplateData{
		Name:       name,
		VerifyLink: verifyLink,
		AppName:    "Website Simple Ecommerce",
	}
	return RenderEmailTemplate("verification", data)
}

func RenderResetPasswordEmail(name, resetLink string) (string, error) {
	data := EmailTemplateData{
		Name:      name,
		ResetLink: resetLink,
		AppName:   "Website Simple Ecommerce",
	}
	return RenderEmailTemplate("reset-password", data)
}
