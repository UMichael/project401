// package main

// import (
// 	"fmt"
// 	"log"
// 	"strings"

// 	"github.com/PuerkitoBio/goquery"
// )

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	var user User
	router := httprouter.New()
	router.GET("/", user.HomePage)
	router.POST("/changePass", user.ChangePassword)
	router.GET("/changePass", user.ChangePasswordGet)
	router.GET("/login", LoginGet)
	router.POST("/login", user.LoginPost)
	router.ServeFiles("/bootstrap4/*filepath", http.Dir("./templates/bootstrap4"))
	router.ServeFiles("/js/*filepath", http.Dir("./templates/js"))
	http.ListenAndServe(":8080", router)
}

// var body = strings.NewReader(`
//         <html>
//         <body>
// 		<table>
// 		<tbody>
//         <tr>
//         <td>Row 1, Content 1</td>
//         <td>Row 1, Content 2</td>
//         <td>Row 1, Content 3</td>
//         <td>Row 1, Content 4</td>
// 		</tr>
// 		</tbody>
// <div>
// 		<tbody id="mic">
//         <tr>
//         <td>Row 2, Content 1</td>
//         <td>Row 2, Content 2</td>
//         <td>Row 2, Content 3</td>
//         <td>Row 2, Content 4</td>
// 		</tr>
// 		<tr>
//         <td>Row 3, Content 1</td>
//         <td>Row 3, Content 2</td>
//         <td>Row 3, Content 3</td>
//         <td>Row 3, Content 4</td>
// 		</tr>
// 		</tbody>
// </div>
//         </table>
//         </body>
//         </html>`)

// func main() {

// 	doc, err := goquery.NewDocumentFromReader(body)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	doc.Find("table").Each(func(i int, tableHtml *goquery.Selection) {
// 		tableHtml.Find("tbody#mic").Each(func(j int, tbodyHtml *goquery.Selection) {
// 			tbodyHtml.Find("tr").Each(func(j int, trHtml *goquery.Selection) {
// 				trHtml.Find("td").Each(func(j int, tdHtml *goquery.Selection) {
// 					fmt.Println(j, strings.TrimSpace(tdHtml.Text()))
// 				})
// 				fmt.Println("done")
// 			})
// 		})
// 	})

// 	// z := html.NewTokenizer(body)
// 	// content := []string{}
// 	// tt := z.Next()
// 	// for {
// 	// 	if tt == html.EndTagToken && z.Token().Data == "tbody" {
// 	// 		break
// 	// 	}
// 	// 	tt = z.Next()
// 	// }
// 	// // While have not hit the </html> tag
// 	// for z.Token().Data != "html" {
// 	// 	tt = z.Next()

// 	// 	// for tt != html.EndTagToken {
// 	// 	// 	tt = z.Next()
// 	// 	// 	t := z.Token()
// 	// 	// 	//fmt.Println(t.Data, "found")
// 	// 	// 	//break
// 	// 	// 	if t.Data == "tbody" {
// 	// 	// 		fmt.Println("efm")
// 	// 	// 		break
// 	// 	// 	}
// 	// 	// }

// 	// 	if tt == html.StartTagToken {
// 	// 		t := z.Token()

// 	// 		if t.Data == "tr" {
// 	// 			innerTag := z.Next()

// 	// 			fmt.Println(innerTag)
// 	// 			if innerTag == html.TextToken {
// 	// 				z.Next()
// 	// 				//fmt.Println(1, inner)
// 	// 				fmt.Println(z.Token().Data)
// 	// 				inner := z.Next()
// 	// 				//	fmt.Println(z.Token().Data)
// 	// 				for inner == html.TextTokenl {
// 	// 					i := 0
// 	// 					i < 4
// 	// 					inner = z.Next()
// 	// 					text := (string)(z.Text())
// 	// 					fmt.Println((text))
// 	// 					i++
// 	// 				}

// 	// 				//t := strings.TrimSpace(text)
// 	// 				//content = append(content, t)
// 	// 			}
// 	// 		}
// 	// 	}
// 	// }
// 	// // Print to check the slice's content
// 	// fmt.Println(content)
// }
