const input = document.getElementById("search-input")
const searchIcon = document.getElementById("search-icon")

input.addEventListener('keydown', function (event) {
    if (event.keyCode === 13 && event.target.value) {
		search({q: event.target.value})
    }
})

searchIcon.addEventListener('click', function (e) {
	if (input.value) {
		search({q: input.value})
	}
})

function search({q, p}) {
	if (p) {
		window.location.href = window.origin + "/search?q=" + encodeURIComponent(q) + "&page=" + p
	} else {
		window.location.href = window.origin + "/search?q=" + encodeURIComponent(q)
	}
}

/**
 * 高亮关键字
 * @param node 节点
 * @param pattern 匹配的正则表达式
 * @param index - 可选。本项目中特定的需求，表示第几组关键词
 * @returns exposeCount - 露出次数
 */
function highlightKeyword(node, pattern) {
	let exposeCount = 0;
	if (node.nodeType === 3) {
		const matchResult = node.data.match(pattern);
		if (matchResult) {
			const highlightEl = document.createElement('span');
			highlightEl.dataset.highlight = 'yes';
			highlightEl.style.color = '#ea4335';
			highlightEl.dataset.highlightMatch = matchResult[0];
			// (index == null) || highlightEl.dataset.highlightIndex = index;
			const matchNode = node.splitText(matchResult.index);
			matchNode.splitText(matchResult[0].length);
			const highlightTextNode = document.createTextNode(matchNode.data);
			highlightEl.appendChild(highlightTextNode);
			matchNode.parentNode.replaceChild(highlightEl, matchNode);
			exposeCount++;
		}
	}
	// 具体条件自己加，这里是基础条件
	else if ((node.nodeType === 1)  && !(/script|style/.test(node.tagName.toLowerCase()))) {
		if (node.dataset.highlight === 'yes') {
			// if (index == null) {
			// 	return;
			// }
			// if (node.dataset.highlightIndex === index.toString()) {
			// 	return;
			// }
			return;
		}
		let childNodes = node.childNodes;
		for (let i = 0; i < childNodes.length; i++) {
			highlightKeyword(childNodes[i], pattern);
		}
	}
	// return exposeCount;
}

/**
 * @param {String | Array} keywords - 要高亮的关键词或关键词数组
 * @returns {Array}
 */
function hanldeKeyword(keywords) {
	let wordMatchString = '';
	const words = [].concat(keywords);
	words.forEach(item => {
		let transformString = item.replace(/[.[*?+^$|()/]|\]|\\/g, '\\$&');
		wordMatchString += `|(${transformString})`;
	});
	wordMatchString = wordMatchString.substring(1);
	// 用于再次高亮与关闭的关键字作为一个整体的匹配正则
	const wholePattern = new RegExp(`^${wordMatchString}$`, 'i');
	// 用于第一次高亮的关键字匹配正则
	const pattern = new RegExp(wordMatchString, 'i');
	return [pattern, wholePattern];
}

if (!searchPage || searchPage === '1') {
	fetch('https://api.woc.cool/codeAnswer?q=' + searchKw, {
		method: "GET",
		headers: {},
	}).then(async res => {
		const dataJson = await res.json()
		if (dataJson.code !== 0) {
			return
		}

		const data = dataJson.data

		document.getElementById("code-block-container").innerHTML = `
		<div class="code-block-container">
			<div class="flex-row-center" style="margin-bottom: 10px;">
				<div class="code-block-title">代码智能搜索结果</div>
			</div>
			<div class="answer-content-container">
				<div class="labels-container">
					${data.tags.map((e) => {
				return `
						<div class="flex-row-center" style="margin-right: 20px;">
							<span class="label-icon" style="background-color: rgb(40, 144, 238);"></span>
							${e}
						</div>
						`
			}).join(" ")}
				</div>
				<div class="flex-row" style="margin-left: 10px;">
					<div class="code-list-container">
						<div class="code-description" style="margin-bottom: 15px; display: inline; line-height: 1.5;">
							${data.ans}
						</div>
					</div>
				</div>
			</div>
			<div class="original-link flex-row-center">
				来源：
				<a href="${data.url}" target="_blank">${data.title}</a>
			</div>
			<div class="flex-row-center" style="max-width: 590px; margin-top: 16px; padding-right: 10px;">
				<div style="width: 100%; height: 1px; background: rgb(235, 235, 235);"></div>
				<div style="margin-left: 12px;"></div>
				<div class="search-results-feedback-text"><span>提供反馈</span></div>
			</div>
		</div>
		`
	})
}

const list = document.getElementsByClassName('snippet')
for (let i = 0; i < searchTags.length; i++) {
	const patterns = hanldeKeyword(searchTags[i]);
	for (let j = 0; j < list.length; j++) {
		highlightKeyword(list[j], patterns[0]);
	}
}

const pageList = document.getElementById("page-list")
for (let i = 1; i < 7; i++) {
	let select = ''
	if ((searchPage == i) || (!searchPage && i === 1)) {
		select = 'selected'
	}
	pageList.innerHTML += `<div class="flex-double-center turn-page-num-wrap ${select}" data-page="${i}">${i}</div>`
}
pageList.childNodes.forEach((node) => {
	node.addEventListener('click', function (e) {
		search({q: searchKw, p: e.target.dataset.page})
	})
})