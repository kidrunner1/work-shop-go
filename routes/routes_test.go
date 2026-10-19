package routes_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"workshop-go/controller"
	"workshop-go/routes"

	"github.com/gofiber/fiber/v2"
)

func TestBasicAuth(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: controller.ErrorHandler})
	routes.Setup(app)

	for _, endpoint := range []struct{ method, path string }{
		{"GET", "/api/v1/fact/5"},
		{"GET", "/api/v3/krit?tax_id=Hello"},
		{"POST", "/api/v1/register"},
	} {
		for _, auth := range []string{"missing", "wrong password", "wrong username"} {
			t.Run(endpoint.path+"/"+auth, func(t *testing.T) {
				req := httptest.NewRequest(endpoint.method, endpoint.path, nil)
				if auth == "wrong password" {
					req.SetBasicAuth("gofiber", "wrong")
				} else if auth == "wrong username" {
					req.SetBasicAuth("other", "21022566")
				}
				res, err := app.Test(req)
				if err != nil {
					t.Fatal(err)
				}
				defer res.Body.Close()
				if res.StatusCode != http.StatusUnauthorized {
					t.Fatalf("status = %d, want 401", res.StatusCode)
				}
				var body map[string]string
				if err := json.NewDecoder(res.Body).Decode(&body); err != nil || body["message"] == "" {
					t.Fatal("401 must return a JSON error message")
				}
			})
		}
	}
}

func TestServerErrors(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: controller.ErrorHandler})
	routes.Setup(app)
	// เส้นทางจำลองมีเฉพาะในชุดทดสอบ ไม่เปิดให้เรียกใน server จริง
	app.Get("/test/error", func(ctx *fiber.Ctx) error {
		return io.ErrUnexpectedEOF
	})
	app.Get("/test/panic", func(ctx *fiber.Ctx) error {
		panic("private internal details")
	})

	for _, path := range []string{"/test/error", "/test/panic"} {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest("GET", path, nil)
			req.SetBasicAuth("gofiber", "21022566")
			res, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			var body map[string]string
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != 500 || body["message"] != "เกิดข้อผิดพลาดภายในระบบ กรุณาลองใหม่อีกครั้ง" {
				t.Fatalf("unexpected response: %d %v", res.StatusCode, body)
			}
		})
	}
}

func TestRegistrationFieldMessage(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: controller.ErrorHandler})
	routes.Setup(app)
	req := httptest.NewRequest("POST", "/api/v1/register", strings.NewReader(`{"username":"student01","email":"bad-email","password":"password123","confirm_password":"other","phone":"0812345678","business_type":"ร้านค้าออนไลน์","website_name":"student-shop"}`))
	req.SetBasicAuth("gofiber", "21022566")
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var body struct {
		Errors map[string]string `json:"errors"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 400 || body.Errors["email"] == "" || body.Errors["confirm_password"] == "" || len(body.Errors) != 2 {
		t.Fatalf("expected email and confirmation errors, got %d %v", res.StatusCode, body.Errors)
	}
}

func TestAPI(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: controller.ErrorHandler})
	routes.Setup(app)

	valid := `{"username":"student01","email":"student@example.com","password":"password123","confirm_password":"password123","phone":"0812345678","business_type":"ร้านค้าออนไลน์","website_name":"student-shop"}`
	tests := []struct {
		name, method, path, body string
		status                   int
		want                     string
	}{
		{"factorial 5", "GET", "/api/v1/fact/5", "", 200, `"factorial":120`},
		{"factorial 0", "GET", "/api/v1/fact/0", "", 200, `"factorial":1`},
		{"factorial 20", "GET", "/api/v1/fact/20", "", 200, `"factorial":2432902008176640000`},
		{"negative", "GET", "/api/v1/fact/-1", "", 400, ""},
		{"overflow", "GET", "/api/v1/fact/21", "", 400, ""},
		{"decimal", "GET", "/api/v1/fact/1.5", "", 400, ""},
		{"not a number", "GET", "/api/v1/fact/abc", "", 400, ""},
		{"ascii letters", "GET", "/api/v3/krit?tax_id=Hello", "", 200, `"ascii":"72 101 108 108 111"`},
		{"ascii mixed", "GET", "/api/v3/krit?tax_id=A12", "", 200, `"ascii":"65 49 50"`},
		{"missing query", "GET", "/api/v3/krit", "", 400, ""},
		{"empty query", "GET", "/api/v3/krit?tax_id=", "", 400, ""},
		{"non ascii", "GET", "/api/v3/krit?tax_id=%E0%B8%81", "", 400, ""},
		{"valid registration", "POST", "/api/v1/register", valid, 200, "ข้อมูลสมัครสมาชิกถูกต้อง"},
		{"missing fields", "POST", "/api/v1/register", `{}`, 400, "ใส่ข้อมูลผิดพลาด"},
		{"malformed json", "POST", "/api/v1/register", `{"username":`, 400, "ใส่ข้อมูลผิดพลาด"},
		{"wrong field type", "POST", "/api/v1/register", `{"username":123}`, 400, "ใส่ข้อมูลผิดพลาด"},
		{"short username", "POST", "/api/v1/register", strings.Replace(valid, "student01", "ab", 1), 400, "ใส่ข้อมูลผิดพลาด"},
		{"long username", "POST", "/api/v1/register", strings.Replace(valid, "student01", strings.Repeat("a", 21), 1), 400, "ใส่ข้อมูลผิดพลาด"},
		{"username symbols", "POST", "/api/v1/register", strings.Replace(valid, "student01", "student!", 1), 400, "ใส่ข้อมูลผิดพลาด"},
		{"invalid email", "POST", "/api/v1/register", strings.Replace(valid, "student@example.com", "bad-email", 1), 400, "ใส่ข้อมูลผิดพลาด"},
		{"short password", "POST", "/api/v1/register", strings.ReplaceAll(valid, "password123", "short"), 400, "ใส่ข้อมูลผิดพลาด"},
		{"long password", "POST", "/api/v1/register", strings.ReplaceAll(valid, "password123", strings.Repeat("a", 73)), 400, "ใส่ข้อมูลผิดพลาด"},
		{"password mismatch", "POST", "/api/v1/register", strings.Replace(valid, `"confirm_password":"password123"`, `"confirm_password":"different"`, 1), 400, "ใส่ข้อมูลผิดพลาด"},
		{"missing confirmation", "POST", "/api/v1/register", strings.Replace(valid, `,"confirm_password":"password123"`, "", 1), 400, "ใส่ข้อมูลผิดพลาด"},
		{"optional line id", "POST", "/api/v1/register", strings.Replace(valid, `"phone":`, `"line_id":"student_line","phone":`, 1), 200, `"line_id":"student_line"`},
		{"phone preserves leading zero", "POST", "/api/v1/register", valid, 200, `"phone":"0812345678"`},
		{"missing phone", "POST", "/api/v1/register", strings.Replace(valid, `,"phone":"0812345678"`, "", 1), 400, `"phone":"กรุณากรอกหมายเลขโทรศัพท์เป็นตัวเลขเท่านั้น"`},
		{"invalid phone", "POST", "/api/v1/register", strings.Replace(valid, "0812345678", "081abc5678", 1), 400, `"phone":`},
		{"decimal phone", "POST", "/api/v1/register", strings.Replace(valid, "0812345678", "081.2345678", 1), 400, `"phone":`},
		{"signed phone", "POST", "/api/v1/register", strings.Replace(valid, "0812345678", "-0812345678", 1), 400, `"phone":`},
		{"missing business type", "POST", "/api/v1/register", strings.Replace(valid, `,"business_type":"ร้านค้าออนไลน์"`, "", 1), 400, `"business_type":"กรุณาเลือกประเภทธุรกิจ"`},
		{"missing website name", "POST", "/api/v1/register", strings.Replace(valid, `,"website_name":"student-shop"`, "", 1), 400, `"website_name":`},
		{"short website name", "POST", "/api/v1/register", strings.Replace(valid, "student-shop", "a", 1), 400, `"website_name":`},
		{"minimum website name", "POST", "/api/v1/register", strings.Replace(valid, "student-shop", "a1", 1), 200, `"website_name":"a1"`},
		{"maximum website name", "POST", "/api/v1/register", strings.Replace(valid, "student-shop", strings.Repeat("a", 30), 1), 200, "ข้อมูลสมัครสมาชิกถูกต้อง"},
		{"long website name", "POST", "/api/v1/register", strings.Replace(valid, "student-shop", strings.Repeat("a", 31), 1), 400, `"website_name":`},
		{"uppercase website name", "POST", "/api/v1/register", strings.Replace(valid, "student-shop", "Student", 1), 400, `"website_name":`},
		{"website leading hyphen", "POST", "/api/v1/register", strings.Replace(valid, "student-shop", "-student", 1), 400, `"website_name":`},
		{"website trailing hyphen", "POST", "/api/v1/register", strings.Replace(valid, "student-shop", "student-", 1), 400, `"website_name":`},
		{"website underscore", "POST", "/api/v1/register", strings.Replace(valid, "student-shop", "student_shop", 1), 400, `"website_name":`},
		{"website thai characters", "POST", "/api/v1/register", strings.Replace(valid, "student-shop", "ร้านค้า", 1), 400, `"website_name":`},
		{"website full domain", "POST", "/api/v1/register", strings.Replace(valid, "student-shop", "student.sogoodweb.com", 1), 400, `"website_name":`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.SetBasicAuth("gofiber", "21022566")
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != tt.status || !strings.Contains(string(body), tt.want) {
				t.Fatalf("status = %d, body = %s; want status %d containing %q", res.StatusCode, body, tt.status, tt.want)
			}
			if !json.Valid(body) {
				t.Fatalf("response is not JSON: %s", body)
			}
			if tt.path == "/api/v1/register" && (strings.Contains(string(body), "password123") || strings.Contains(string(body), "different")) {
				t.Fatal("response must not contain a password")
			}
		})
	}
}
