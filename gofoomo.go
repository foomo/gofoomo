// Integrate golang with the foomo php framework. Think of gophers riding elephants or
// maybe also think of gophers pulling toy elephants.
//
//
// Example:
//
// 	f := gofoomo.NewFoomo("/Users/jan/vagrant/schild/www/schild", "test")
// 	p := proxy.NewProxy(f, "http://schild-local-test.bestbytes.net")
//	// the static files handler will keep requests to static files away from apache
// 	p.AddHandler(handler.NewStaticFiles(f))
// 	http.ListenAndServe(":8080", p)
//
package gofoomo
