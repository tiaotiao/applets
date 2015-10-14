package main

import (
	"github.com/tiaotiao/web"
)

func IndexPage(c *web.Context) interface{} {
	return indexPage
}

var indexPage = `
<html>
	<head>
		<title></title>
	</head>
	<body>
		<h2>Submit your code</h2>

		<p>Run any code as you wish, use a global variable named 'result' to return a result back.</p>

		<form enctype="text/plain" id="codeform" method="get">
			<p>
				<strong>Select your language:</strong>
				<input name="lang" type="radio" value="js" checked="checked" />javascript
				<input name="lang" type="radio" value="go" />golang 
				<input name="lang" type="radio" value="python" />python
				<input name="lang" type="radio" value="lua" />lua
			</p>

			<p><textarea id="code" cols="100" name="code" rows="30"></textarea></p>
			<p><input type="submit" value="     Submit     " /></p>
		</form>

		<p><b>Status:</b> <span id="A1"></span></p>
		<p><b>Status text:</b> <span id="A2"></span></p>
		<p><b>Response:</b> <span id="A3"></span></p>

		<p>&nbsp;</p>
	</body>

	<script>
	function request() {
	  var xhttp = new XMLHttpRequest();
	  xhttp.onreadystatechange = function() {
	      document.getElementById('A1').innerHTML = xhttp.status;
	      document.getElementById('A2').innerHTML = xhttp.statusText;
	      document.getElementById('A3').innerHTML = xhttp.responseText;
	  }

	  var lang = document.querySelector('input[name="lang"]:checked').value;
	  var code = document.getElementById('code').value;
	  var params = "lang=" + lang + "&code=" + encodeURIComponent(code);

	  xhttp.open("POST", "/api/run");
	  xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
	  xhttp.send(params);
	}

	window.addEventListener("load", function() {
		var form = document.getElementById("codeform");
		// to takeover its submit event.
		form.addEventListener("submit", function (event) {
			event.preventDefault();
			request();
		});
	})
	</script>
</html>
`
