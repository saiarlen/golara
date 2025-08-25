package main

import (
	"fmt"
	"html/template"
	"os"
)

type Documentation struct {
	Title       string
	Version     string
	Description string
	Sections    []Section
}

type Section struct {
	ID      string
	Title   string
	Content string
}

func main() {
	docs := Documentation{
		Title:       "Golara Framework",
		Version:     "1.0.0",
		Description: "Complete Laravel-inspired Go MVC framework",
		Sections: []Section{
			{
				ID:    "installation",
				Title: "Installation",
				Content: `
<p>Get started with Golara in minutes:</p>
<pre><code class="language-bash"># Clone the framework
git clone https://github.com/your-repo/golara
cd golara

# Initialize your project
go run main.go -subcommand init

# Install dependencies
go mod tidy

# Build and run
make build
./bin/golara</code></pre>`,
			},
			{
				ID:    "mvc",
				Title: "Complete MVC Architecture",
				Content: `
<div class="mvc-section">
    <h3>🎨 Models, Views & Controllers</h3>
    <div class="feature-grid">
        <div class="feature-card">
            <h4>M - Models</h4>
            <p>Eloquent-style models with GORM</p>
            <pre><code>./bin/golara -subcommand make:model Product</code></pre>
        </div>
        <div class="feature-card">
            <h4>V - Views</h4>
            <p>Laravel Blade-like templates</p>
            <pre><code>./bin/golara -subcommand make:view products/index</code></pre>
        </div>
        <div class="feature-card">
            <h4>C - Controllers</h4>
            <p>API & Web resource controllers</p>
            <pre><code>./bin/golara -subcommand make:controller Product --web</code></pre>
        </div>
    </div>
</div>`,
			},
			{
				ID:    "views",
				Title: "Template Engine",
				Content: `
<h3>Laravel Blade-like Templates</h3>
<pre><code class="language-html"><!-- resources/views/products/index.html -->
{{define "content"}}
<h1>{{.title}}</h1>
{{range .products}}
    <div class="product">
        <h2>{{.Name}}</h2>
        <p>Price: ${{.Price}}</p>
        <a href="{{route "products.show"}} {{.ID}}">View</a>
    </div>
{{end}}
{{end}}</code></pre>

<h3>Web Controller</h3>
<pre><code class="language-go">func (ctrl *ProductController) Index(c *fiber.Ctx) error {
    var products []models.Product
    // ... fetch data
    
    return view.View(c, "products/index", view.ViewData{
        "products": products,
        "title": "All Products",
    }, "layouts/app")
}</code></pre>

<h3>Helper Functions</h3>
<ul>
    <li><code>{{asset "css/app.css"}}</code> - Asset URLs</li>
    <li><code>{{route "home"}}</code> - Named routes</li>
    <li><code>{{csrf}}</code> - CSRF token</li>
</ul>`,
			},
			{
				ID:    "features",
				Title: "Framework Features",
				Content: `
<div class="feature-grid">
    <div class="feature-card">
        <h3>🗄️ Laravel-style ORM</h3>
        <p>Fluent query builder with relationships and migrations</p>
    </div>
    <div class="feature-card">
        <h3>⚡ Redis Integration</h3>
        <p>Built-in caching and queue management</p>
    </div>
    <div class="feature-card">
        <h3>🎨 Template Engine</h3>
        <p>Blade-like templates with layouts and helpers</p>
    </div>
    <div class="feature-card">
        <h3>📡 Event System</h3>
        <p>Observer pattern for decoupled architecture</p>
    </div>
    <div class="feature-card">
        <h3>🏗️ Service Container</h3>
        <p>Dependency injection and service providers</p>
    </div>
    <div class="feature-card">
        <h3>🛡️ Security Features</h3>
        <p>JWT auth, rate limiting, CORS, validation</p>
    </div>
</div>`,
			},
		},
	}

	generateHTML(docs)
	generateCSS()
	generateJS()
	
	fmt.Println("✅ Documentation generated successfully!")
	fmt.Println("📖 Open docs/html/index.html in your browser")
}

func generateHTML(docs Documentation) {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} Documentation</title>
    <link rel="stylesheet" href="../assets/css/docs.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.24.1/themes/prism.min.css">
</head>
<body>
    <nav class="sidebar">
        <div class="logo">
            <h1>{{.Title}}</h1>
            <span class="version">v{{.Version}}</span>
        </div>
        <ul class="nav-menu">
            {{range .Sections}}
            <li><a href="#{{.ID}}">{{.Title}}</a></li>
            {{end}}
        </ul>
    </nav>
    
    <main class="content">
        <header class="hero">
            <h1>{{.Title}}</h1>
            <p class="lead">{{.Description}}</p>
            <div class="hero-buttons">
                <a href="#installation" class="btn btn-primary">Get Started</a>
                <a href="https://github.com/your-repo/golara" class="btn btn-secondary">GitHub</a>
            </div>
        </header>
        
        {{range .Sections}}
        <section id="{{.ID}}" class="doc-section">
            <h2>{{.Title}}</h2>
            <div class="section-content">
                {{.Content}}
            </div>
        </section>
        {{end}}
    </main>
    
    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.24.1/components/prism-core.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.24.1/plugins/autoloader/prism-autoloader.min.js"></script>
    <script src="../assets/js/docs.js"></script>
</body>
</html>`

	t, err := template.New("docs").Parse(tmpl)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("docs/html/index.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = t.Execute(file, docs)
	if err != nil {
		panic(err)
	}
}

func generateCSS() {
	css := `/* Golara Documentation Styles */
:root {
    --primary-color: #ff6b35;
    --secondary-color: #2c3e50;
    --text-color: #333;
    --bg-color: #f8f9fa;
    --sidebar-width: 280px;
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    line-height: 1.6;
    color: var(--text-color);
    background: var(--bg-color);
}

.sidebar {
    position: fixed;
    top: 0;
    left: 0;
    width: var(--sidebar-width);
    height: 100vh;
    background: white;
    border-right: 1px solid #e1e8ed;
    padding: 2rem 0;
    overflow-y: auto;
}

.logo {
    padding: 0 2rem 2rem;
    border-bottom: 1px solid #e1e8ed;
    margin-bottom: 2rem;
}

.logo h1 {
    color: var(--primary-color);
    font-size: 1.5rem;
    margin-bottom: 0.5rem;
}

.version {
    background: var(--primary-color);
    color: white;
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    font-size: 0.75rem;
}

.nav-menu {
    list-style: none;
}

.nav-menu a {
    display: block;
    padding: 0.75rem 2rem;
    color: var(--text-color);
    text-decoration: none;
    transition: all 0.2s;
}

.nav-menu a:hover {
    background: var(--bg-color);
    color: var(--primary-color);
}

.content {
    margin-left: var(--sidebar-width);
    padding: 2rem;
    max-width: 800px;
}

.hero {
    text-align: center;
    padding: 4rem 0;
    margin-bottom: 4rem;
}

.hero h1 {
    font-size: 3rem;
    color: var(--secondary-color);
    margin-bottom: 1rem;
}

.lead {
    font-size: 1.25rem;
    color: #666;
    margin-bottom: 2rem;
}

.hero-buttons {
    display: flex;
    gap: 1rem;
    justify-content: center;
}

.btn {
    padding: 0.75rem 2rem;
    border-radius: 6px;
    text-decoration: none;
    font-weight: 500;
    transition: all 0.2s;
}

.btn-primary {
    background: var(--primary-color);
    color: white;
}

.btn-primary:hover {
    background: #e55a2b;
}

.btn-secondary {
    background: transparent;
    color: var(--secondary-color);
    border: 2px solid var(--secondary-color);
}

.btn-secondary:hover {
    background: var(--secondary-color);
    color: white;
}

.doc-section {
    margin-bottom: 4rem;
}

.doc-section h2 {
    color: var(--secondary-color);
    font-size: 2rem;
    margin-bottom: 1.5rem;
    padding-bottom: 0.5rem;
    border-bottom: 2px solid var(--primary-color);
}

.feature-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1.5rem;
    margin: 2rem 0;
}

.feature-card {
    background: white;
    padding: 1.5rem;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.feature-card h3, .feature-card h4 {
    color: var(--primary-color);
    margin-bottom: 0.5rem;
}

pre {
    background: #2d3748;
    color: #e2e8f0;
    padding: 1.5rem;
    border-radius: 6px;
    overflow-x: auto;
    margin: 1rem 0;
}

code {
    font-family: 'Monaco', 'Menlo', monospace;
    font-size: 0.9rem;
}

@media (max-width: 768px) {
    .sidebar {
        transform: translateX(-100%);
    }
    
    .content {
        margin-left: 0;
        padding: 1rem;
    }
    
    .hero h1 {
        font-size: 2rem;
    }
    
    .hero-buttons {
        flex-direction: column;
        align-items: center;
    }
}`

	err := os.WriteFile("docs/assets/css/docs.css", []byte(css), 0644)
	if err != nil {
		panic(err)
	}
}

func generateJS() {
	js := `// Golara Documentation JavaScript
document.addEventListener('DOMContentLoaded', function() {
    // Smooth scrolling for navigation links
    const navLinks = document.querySelectorAll('.nav-menu a');
    
    navLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            const targetId = this.getAttribute('href').substring(1);
            const targetElement = document.getElementById(targetId);
            
            if (targetElement) {
                targetElement.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        });
    });
    
    // Highlight active section in navigation
    const sections = document.querySelectorAll('.doc-section');
    const navItems = document.querySelectorAll('.nav-menu a');
    
    function highlightActiveSection() {
        let current = '';
        
        sections.forEach(section => {
            const sectionTop = section.offsetTop;
            const sectionHeight = section.clientHeight;
            
            if (window.pageYOffset >= sectionTop - 200) {
                current = section.getAttribute('id');
            }
        });
        
        navItems.forEach(item => {
            item.classList.remove('active');
            if (item.getAttribute('href') === '#' + current) {
                item.classList.add('active');
            }
        });
    }
    
    window.addEventListener('scroll', highlightActiveSection);
    highlightActiveSection();
});`

	err := os.WriteFile("docs/assets/js/docs.js", []byte(js), 0644)
	if err != nil {
		panic(err)
	}
}