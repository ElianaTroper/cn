
// TODO: Create a post function, which posts as needed

function post(url, content) {
	const xh = new XMLHttpRequest();
	xh.open('POST', url);
	// XXX: Only posts JSON - for now this is all we want.
	xh.setRequestHeader('content-type', 'application/json');
	xh.send(JSON.stringify(content))
}
