// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.906
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Index() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<script>\n            // Forward scroll events to the right column when scrolling anywhere on the page\n            document.addEventListener('wheel', function(e) {\n                const scrollableContent = document.getElementById('scrollable-content');\n                const isDesktop = window.matchMedia('(min-width: 768px)').matches;\n\n                if (scrollableContent && !scrollableContent.contains(e.target) && isDesktop) {\n                    // Prevent default scrolling behavior\n                    e.preventDefault();\n                    \n                    // Forward the scroll to the right column\n                    scrollableContent.scrollTop += e.deltaY;\n                }\n            }, { passive: false });\n        </script> <main class=\"flex flex-col md:flex-row h-full md:overflow-hidden\"><div class=\"flex flex-col gap-8 md:gap-0 items-center justify-evenly p-8 md:p-4 bg-neutral text-neutral-content md:sticky top-0 md:w-1/2 md:max-h-full md:h-full\"><div class=\"border-b border-primary px-0 py-4 md:p-4\"><h1 class=\"text-4xl py-1\">Oscar Korpi</h1><div><h2 class=\"text-2xl py-1\">M.Sc. Student in Computer Science & Engineering at LTH</h2><p class=\"py-1\">Specializing in software, with a focus on cloud, databases and statistics, to build <br>the data-driven applications of the future.</p></div></div><div><a href=\"mailto:contact@korpi.se\">contact@korpi.se</a></div><div class=\"absolute left-0 bottom-0 p-2 text-sm text-neutral-content opacity-70\">Copyright &copy; Oscar Korpi 2025</div></div><article class=\"grid grid-cols-1 items-center p-8 md:px-24 md:w-1/2 overflow-y-scroll divide-y divide-primary\" id=\"scrollable-content\"><section class=\"py-8\"><h2 class=\"text-2xl text-gray-900\">About me</h2><p class=\"py-2\">Hello! I'm Oscar, a student in Computer Science & Engineering at LTH, Sweden.  I'm interested in backend and software development, finance, data science, and statistics. I enjoy building backend applications and exploring the intersection of technology and finance in my free time.<br><br>This website serves as a personal portfolio and blog where I share my journey, projects, and on this website, you'll find my resume, personal projects, and some of my thoughts on various topics. Feel free to reach out if you have any questions or just want to chat! <ul class=\"list-disc pl-6 py-2\"><li><a class=\"link\" href=\"https://www.linkedin.com/in/oscar-korpi-421841234\">Linkedin</a></li><li><a class=\"link\" href=\"https://github.com/o-korpi\">Github</a></li><li><a class=\"link\" href=\"mailto:contact@korpi.se\">Email</a></li></ul></p></section><section class=\"py-8\"><h2 class=\"text-2xl text-gray-900\">Professional experience</h2><ul><li class=\"flex flex-col py-4\"><div class=\"flex flex-col gap-1 md:gap-0 md:flex-row md:justify-between text-xl py-2\"><h3 class=\"text-gray-900\">Software Engineer &ndash; Nordic Civil Engineering</h3><p>2025/02 &ndash; Now</p></div><p>Working with developing internal tooling.</p><ul class=\"join join-horizontal overflow-x-scroll md:overflow-x-auto py-2\"><li class=\"badge badge-neutral join-item bg-[#B125EA] border-[#B125EA]\">Kotlin&trade;</li><li class=\"badge badge-neutral join-item bg-[#9179E4] border-[#9179E4]\">C#</li><li class=\"badge badge-neutral join-item bg-[#512BD4] border-[#512BD4]\">.NET</li><li class=\"badge badge-neutral join-item bg-[#104581] border-[#104581]\">Azure</li><li class=\"badge badge-neutral join-item bg-[#336791] border-[#336791]\">PostgreSQL</li><li class=\"badge badge-neutral join-item bg-[#B125EA] border-[#B125EA]\">Ktor</li></ul></li><li class=\"flex flex-col py-4\"><div class=\"flex flex-col gap-1 md:gap-0 md:flex-row md:justify-between text-xl py-2\"><h3 class=\"text-gray-900\">Teaching Assistant &ndash; Lund University</h3><p>2024/06 &ndash; 2024/12</p></div><p>Teaching Assistant at the Department of Computer Science at LTH, the Faculty of Engineering at Lund University. Mainly worked as a lab assistant, helping students in the courses  Introduction to Programming in Scala and Programming, Second Course (Java). </p><ul class=\"join join-horizontal overflow-x-scroll md:overflow-x-auto py-2\"><li class=\"badge badge-neutral join-item bg-[#f89820] border-[#f89820] text-white\">Java</li><li class=\"badge badge-neutral join-item bg-[#DE3423] border-[#DE3423] text-[#380D09]\">Scala</li></ul></li><li class=\"flex flex-col py-4\"><div class=\"flex flex-col gap-1 md:gap-0 md:flex-row md:justify-between text-xl py-2\"><h3 class=\"text-gray-900\">Software Engineer, Summer Intern &ndash; Nordic Civil Engineering</h3><p>2024/06 &ndash; 2024/08</p></div><p>Worked during the summer on to develop the GoGreen Logistics project together with another student. </p><ul class=\"join join-horizontal overflow-x-scroll md:overflow-x-auto py-2\"><li class=\"badge badge-neutral join-item bg-[#B125EA] border-[#B125EA]\">Kotlin&trade;</li><li class=\"badge badge-neutral join-item bg-[#F0DB4F] border-[#F0DB4F] text-[#323330]\">JavaScript</li><li class=\"badge badge-neutral join-item bg-[#104581] border-[#104581]\">Azure</li><li class=\"badge badge-neutral join-item bg-[#B125EA] border-[#B125EA]\">Ktor</li><li class=\"badge badge-neutral join-item bg-[#5B96D5] border-[#5B96D5]\">HTMX</li></ul></li><li class=\"flex flex-col py-4\"><div class=\"flex flex-col gap-1 md:gap-0 md:flex-row md:justify-between text-xl py-2\"><h3 class=\"text-gray-900\">Teaching Assistant &ndash; Lund University</h3><p>2023/08 &ndash; 2024/03</p></div><p>Teaching Assistant at the Department of Computer Science at LTH, the Faculty of Engineering at Lund University. Worked as a lab assistant, helping students in the courses  Introduction to Programming in Scala and Programming, Second Course (Java). </p><ul class=\"join join-horizontal overflow-x-scroll md:overflow-x-auto py-2\"><li class=\"badge badge-primary join-item bg-[#f89820] border-[#f89820] text-white\">Java</li><li class=\"badge badge-primary join-item bg-[#DE3423] border-[#DE3423] text-[#380D09]\">Scala</li></ul></li></ul></section><section class=\"py-8\"><h2 class=\"text-2xl text-gray-900\">Education</h2><ul><li class=\"flex flex-col py-4\"><div class=\"w-full flex flex-col gap-1 md:gap-0 md:flex-row md:justify-between text-xl py-2\"><h3 class=\"text-gray-900\">M.Sc. in Computer Science & Engineering at LTH</h3><p>2022&ndash;2027</p></div><p>Currently studying, with a planned specialization in Software. Additionally taking a lot of courses in statistics.</p><p>Expected graduation: 2027<br></p><p><br>Notable completed courses:</p><ul class=\"list-disc pl-6 py-2\"><li>Software Development in Teams</li></ul><p>Notable planned master's courses:</p><ul class=\"list-disc pl-6 py-2\"><li>Cloud Computing</li><li>Database Technology</li><li>Applied Machine Learning</li><li>Time Series Analysis</li><li>Monte Carlo-based Statistical Methods</li><li>Stationary and Non-stationary Spectral Analysis</li><li>Statistical Modelling of Extreme Values</li></ul></li><li class=\"flex flex-col py-4\"><div class=\"flex flex-col gap-1 md:gap-0 md:flex-row md:justify-between text-xl py-2\"><h3 class=\"text-gray-900\">Microeconomics (11hp)</h3><p>2025</p></div></li><li class=\"flex flex-col py-4\"><div class=\"flex flex-col gap-1 md:gap-0 md:flex-row md:justify-between text-xl py-2\"><h3 class=\"text-gray-900\">Managerial Economics, Basic Course (7,5hp)</h3><p>2024</p></div></li></ul></section><section class=\"py-8\"><h2 class=\"text-2xl text-gray-900\">Featured personal projects</h2></section><section class=\"py-8\"><h2 class=\"text-2xl text-gray-900\">About this website</h2><p class=\"py-2\">This website was built using Go, Templ, HTMX and surreal.js. HTMX and surreal.js bring interactivity to the website, while Templ handles the templating.<br><br>The blog pages are created by rendering Markdown and converting it to HTML.  Inline LaTeX support is added using some custom parsing and by using MathJax on the frontend.</p></section></article></main>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return nil
		})
		templ_7745c5c3_Err = Base().Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
