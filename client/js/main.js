
import { NODE } from './core.js';
import { runWiki } from './app/wiki.js'
import { runGutenberg } from './app/gutenberg.js'

async function loadApp(node, params) {

	// TODO: Could use return values to function as indicators of completeness? Or maybe something with the node object?
		// Basically, for loading from cache we need to ensure that there are no pending activities that might bork the page

	const app = params.get('app')
	// FUTURE id to hash
	// body = await doCat(hash)
	// FUTURE: Use this generic function: const runFunc = new Function('NODE', 'path', body ) // All of these files must use NODE as the fixed vars, and path if the URL path is useful
	// runFunc(NODE, path)
	switch (app) {
		case 'wiki':
			runWiki(node, params); // TODO: Finish runWiki
			break;
		case 'gutenberg':
			runGutenberg(node, params) // TODO: Finish
			break;
		case 'post':
			runPost(node, params)
		default:
			throw 'Invalid app'
	}
}

function makeLink(id, text) {
	return '<a href="javascript:void(0);" id="' + id + '">' + text + '</a>'
}

function linksList(links) {
	let linksText = ''
	for (const link of links) {
		linksText += '\n<li>'+link+'</li>'
	}
	return `
		<ul>${linksText}
		</ul>
	`;
}

function updateURL(newParams) {
	window.history.pushState( null, '', '?' + newParams.toString())
}

function addLinksFn(node, links, params, toIndex, key, mod) {
	for (const [_key, value] of Object.entries(toIndex)) {
		let element = document.getElementById(params.app + '-' + _key)
		element.onclick = () => {
			params.set(key, mod+_key)
			updateURL(params)
			loadApp(node, params)
		}
	}
}

export function runIndex(node, params, toIndex, key, mod='') {

	// TODO: Add multipage functionality
	// TODO: add a lastIndex check
	
	let links = []
	
	for (const [_key, value] of Object.entries(toIndex)) {
		links.push(makeLink(params.app + '-' + _key, value))
	}
	
	document.body.innerHTML = linksList(links)
	addLinksFn(node, links, params, toIndex, key, mod) // TODO: Hide the page until this is done
	// XXX: Should keep the page loaded, NOT refresh the page. We want to maintain the node during refreshes
}

function main() {
	
	console.log('main: started')

	const urlParams = new URLSearchParams(window.location.search);

	if (urlParams.get('app')) {
		loadApp(NODE, urlParams);
	} else {
		let toIndex = {}
		for (const [_key, value] of Object.entries(NODE.root.app)) {
			toIndex[_key] = value.name
		}
		console.log(toIndex)
		runIndex(NODE, urlParams, toIndex, 'app')
	}
	// FUTURE: Catch errors like invalid pathname, throw alert
}

main()

window.onpopstate = function(event) {
	console.log('Reloaded, but using cache :)')
	main()
	// XXX: Long running events might cause weird behavior here
};

// TODO: Eventually take the click on apps and adjust the URL

// document.addEventListener('DOMContentLoaded', () => { main() });
