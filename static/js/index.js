// const isDev = window.origin !== "https://wox.cool";
const isSmallDevice = window.innerWidth <= 650;
const inputGroupElement = document.getElementById("input-group");
const inputDivElement = document.getElementById("input-div");
const inputElement = document.getElementById("input-box");
const suggestionElement = document.getElementById("suggestion-container");
const closeElement = document.getElementById("close-icon");
const textLabelElement = document.getElementById("text-label");
const refreshElement = document.getElementById("real-time-refresh");
const refreshImgElement = document.getElementById("real-time-refresh-img");
const realTimeConElement = document.getElementById("real-time-container");
const operationWrapElement = document.getElementById("operations-wrap");
const operationContainerElement = document.getElementById("operations-container");
const moreToolsSwitchElement = document.getElementById("more-tools-switch");
const moreToolsContainerElement = document.getElementById("more-tools-container");
const doggoElement = document.getElementById("doggo");
const doggoToggleElement = document.getElementById("toggle-item");
const doggoToggleToolTipElement = document.getElementById("toggle-tooltip");
const doggoToggleToolTipContainerElement = document.getElementById("toggle-tooltip-container");
const toggleTooltipTextElement = document.getElementById("toggle-tooltip-text");
const footerElement = document.getElementById("footer-info");
const privacyEntryElement = document.getElementById("bottom-privacy-entry");
const recordNumberElement = document.getElementById("record-number");
const bottomRightBarElement = document.getElementById("bottom-right-bar");
const themeContainerElement = document.getElementById("footer-setting-container");
const themeLightElement = document.getElementById("footer-theme-mode-light");
const themeLightImgElement = document.getElementById("footer-theme-mode-img-light");
const themeReadingElement = document.getElementById("footer-theme-mode-reading");
const themeReadingImgElement = document.getElementById("footer-theme-mode-img-reading");
const themeDarkElement = document.getElementById("footer-theme-mode-dark");
const themeDarkImgElement = document.getElementById("footer-theme-mode-img-dark");
const mobileSugElement = document.getElementById("m-suggestion-container");
const mobileSugBackElement = document.getElementById("m-backward-container");
const mobileSugInputElement = document.getElementById("m-input-box");
const mobileSugCloseElement = document.getElementById("m-close-container");
const mobileSugItemsElement = document.getElementById("m-item-container");
const fLogoImgElement = document.getElementById("f-logo-img");
const logoEmojiElement = document.getElementById("logo-emoji");
const logoWrapElement = document.getElementById("logo-wrap");
const userSmallAvatarElement = document.getElementById("user-small-avatar");
const userBigAvatarElement = document.getElementById("user-big-avatar");
const emailDetailElement = document.getElementById("email-detail");
const toolsContainerElement = document.getElementById("tools-container");
initThemeMode();
handleExpiredHistoryKeywords();
const emojiArr = [" CatSearch️"];
const randomIndex = Math.floor(Math.random() * emojiArr.length);
const emojiText = emojiArr[randomIndex];
logoEmojiElement.innerHTML = "<img class='logo-arrow' src='../static/img/home_animation_arrow.svg' />" + emojiText;
fLogoImgElement.classList.add("animation");
logoWrapElement.classList.add("animation");
const enableBetaFunction = localStorage.getItem("enableBetaFunction") === "true";
if (enableBetaFunction) {
    doggoToggleElement.style.backgroundColor = "#1973E8";
    doggoElement.checked = true;
    doggoToggleToolTipContainerElement.style.left = "66px";
}

doggoToggleToolTipElement.addEventListener("mouseenter", function () {
    doggoToggleToolTipContainerElement.style.display = "block";
});
doggoToggleToolTipElement.addEventListener("mouseleave", function () {
    doggoToggleToolTipContainerElement.style.display = "none";
});
let isShowAnimation = false;
logoEmojiElement.addEventListener("animationstart", function () {
    isShowAnimation = true;
});
logoEmojiElement.addEventListener("animationend", function () {
    logoWrapElement.classList.remove("animation");
    logoWrapElement.classList.remove("animation-emoji");
    isShowAnimation = false;
});
logoWrapElement.addEventListener("click", function () {
    if (!isShowAnimation) {
        logoWrapElement.classList.add("animation-emoji");
        const randomIndex = Math.floor(Math.random() * emojiArr.length);
        const emojiText = emojiArr[randomIndex];
        logoEmojiElement.innerHTML = "<img class='logo-arrow' src='../static/img/home_animation_arrow.svg' />" + emojiText;
    }
});

let timer;

function debounce(func, duration) {
    clearTimeout(timer);
    timer = setTimeout(() => {
        func && func.apply(this);
    }, duration);
}

let _suggestionIndex = null;
let suggestions = [];

function handleExpiredHistoryKeywords() {
    const _lastTime = localStorage.getItem("lastCheckExpiredTime");
    if (_lastTime) {
        const lastTime = new Date(_lastTime);
        const todayTime = new Date();
        if (lastTime.getFullYear() === todayTime.getFullYear() &&
            lastTime.getMonth() === todayTime.getMonth() &&
            lastTime.getDate() === todayTime.getDate()) {
            return;
        }
    }
    const _keywords = localStorage.getItem("searchKeywordHistory");
    if (!_keywords || JSON.parse(_keywords).length === 0) {
        return;
    }
    let keywords = JSON.parse(_keywords);
    if (typeof keywords[0] === "string") {
        keywords = keywords.map(e => ({keyword: e, createTime: Date.now()}));
    }
    const validKeywords = keywords.filter(e => {
        const pastTime = Date.now() - e.createTime;
        const pastDays = pastTime / 1000 / 60 / 60 / 24;
        return pastDays < 4;
    });
    localStorage.setItem("searchKeywordHistory", JSON.stringify(validKeywords));
    localStorage.setItem("lastCheckExpiredTime", new Date());
}

function getKeywordSuggestions(keyword) {
    const HOST = "https://devapis.hetaoapis.com/";
    const _engine = "search-engine-listing";
    const FULL_HOST = HOST + _engine + "/v1/";
    return fetch(FULL_HOST + "input/suggestions?keywords=" + keyword, {
        method: "GET",
        headers: {
            "x-ht-env": "prod", // isDev ? "dev" : "prod",
            "x-ht-s": "search",
            "Access-Control-Allow-Origin": "*",
        },
    })
}

async function generateSuggestions(searchText) {
    try {
        let sugResult;
        const showKeywordsHistory = localStorage.getItem("showKeywordsHistory") !== "false";
        const _keywords = showKeywordsHistory ? localStorage.getItem("searchKeywordHistory") : [];
        const keywords = _keywords && showKeywordsHistory ? JSON.parse(_keywords) : [];
        if (searchText.trim()) {
            let localSuggestions = keywords.filter(item => item.keyword && item.keyword.indexOf(searchText) === 0).splice(0, 3).map(item => ({
                text: item.keyword,
                type: "local"
            }));
            const sugRes = await getKeywordSuggestions(searchText);
            const data = await sugRes.json();
            const sugData = (data.suggestions || []).filter(e => !localSuggestions.some(item => item.text == e));
            const serverSuggestions = sugData.map(item => ({text: item, type: "server"}));
            sugResult = [...localSuggestions, ...serverSuggestions];
            sugResult = sugResult.length > 8 ? sugResult.splice(0, 8) : sugResult;
        } else {
            sugResult = keywords.splice(0, 8).map(item => ({text: item.keyword, type: "local"}));
        }
        sugResult = sugResult.filter(e => typeof e.text !== "object");
        sugResult.forEach(e => {
            e.text = e.text.replace(/'/g, "‘").replace(/"/g, "“").replace(/</g, " ").replace(/>/g, " ");
        });
        return sugResult;
    } catch (error) {
        isDev && console.log('error :>> ', error)
    }
}

function setHistoryKeywordsWithSearch(keyword) {
    if (keyword.trim()) {
        const showKeywordsHistory = localStorage.getItem("showKeywordsHistory") !== "false";
        if (!showKeywordsHistory) {
            return;
        }
        const _keywords = localStorage.getItem("searchKeywordHistory");
        let keywords = _keywords ? JSON.parse(_keywords) : [];
        keywords.unshift({keyword: keyword.trim(), createTime: Date.now()});
        const _keywordsArr = [];
        keywords = keywords.filter(e => {
            if (!_keywordsArr.includes(e.keyword)) {
                _keywordsArr.push(e.keyword);
                return true;
            } else {
                return false;
            }
        });
        localStorage.setItem("searchKeywordHistory", JSON.stringify(keywords));
    }
}

function handleDeleteHistoryKeywords(keyword, sugs) {
    let _keywords = localStorage.getItem("searchKeywordHistory");
    _keywords = _keywords ? JSON.parse(_keywords) : [];
    const keywords = _keywords.filter(item => item.keyword !== keyword);
    localStorage.setItem("searchKeywordHistory", JSON.stringify(keywords));
    let _suggestions = JSON.parse(JSON.stringify(sugs));
    _suggestions = _suggestions.filter(item => item.text !== keyword);
    return _suggestions || [];
}

function renderSuggestionsWithData(type, sugs) {
    if (sugs.length !== 0) {
        if (type === "mobile") {
            mobileSugItemsElement.innerHTML = "";
            sugs.forEach(item => {
                if (item.type === "server") {
                    mobileSugItemsElement.innerHTML += "<div id='" + (item.type + item.text) + "'" + "onclick='searchWithValue(" + ('"sugLink"') + "," + ('"mobile"') + "," + ('"' + item.text + '"') + ")' class='m-suggestion-item'>" +
                        "<img class='m-search-icon' src='https://cdn.fsofso.com/static/assets/keyword_search.svg'/>" +
                        "<div class='m-suggestion-text'>" + item.text + "</div>" +
                        "<img class='m-navigate-icon' src='https://cdn.fsofso.com/static/assets/icon_arrow8.svg'/>" +
                        "</div>";
                } else {
                    mobileSugItemsElement.innerHTML += "<div id='" + (item.type + item.text) + "'" + "onclick='searchWithValue(" + ('"sugLink"') + "," + ('"mobile"') + "," + ('"' + item.text + '"') + ")' class='m-suggestion-item'>" +
                        "<img class='m-search-icon' src='https://cdn.fsofso.com/static/assets/history_keyword_icon.svg'/>" +
                        "<div class='m-suggestion-text'>" + item.text + "</div>" +
                        "<img onclick='clickDeleteHistoryKeyword(" + ('"' + item.text + '"') + "," + ('"mobile"') + ")' class='m-navigate-icon' src='https://cdn.fsofso.com/static/assets/history_keyword_delete.svg'/>" +
                        "</div>";
                }
            });
        } else {
            _suggestionIndex = null;
            suggestionElement.innerHTML = "<div class='input-suggestion-divider'></div>";
            sugs.forEach(item => {
                if (item.type === "server") {
                    suggestionElement.innerHTML += "<div id='" + (item.type + item.text) + "'" + "onclick='searchWithValue(" + ('"sugLink"') + "," + ('"pc"') + "," + ('"' + item.text + '"') + ")' class='suggestion-item'>" +
                        "<div class='suggestion-content-wrap'>" +
                        "<img class='suggestion-search-icon' src='https://cdn.fsofso.com/static/assets/search_sug_searchicon.svg'/>" +
                        "<div class='suggestion-value'>" + item.text + "</div>" +
                        "</div>" +
                        "</div>";
                } else {
                    suggestionElement.innerHTML += "<div id='" + (item.type + item.text) + "'" + "onclick='searchWithValue(" + ('"sugLink"') + "," + ('"pc"') + "," + ('"' + item.text + '"') + ")' class='suggestion-item'>" +
                        "<div class='suggestion-content-wrap'>" +
                        "<img class='suggestion-search-icon' src='https://cdn.fsofso.com/static/assets/search_sug_clock.svg'/>" +
                        "<div class='suggestion-value'>" + item.text + "</div>" +
                        "</div>" +
                        "<img id='" + (item.type + item.text + "delete-logo") + "'" + "onclick='clickDeleteHistoryKeyword(" + ('"' + item.text + '"') + "," + ('"pc"') + ")' class='suggestion-delete-icon' src='https://cdn.fsofso.com/static/assets/search_sug_close.svg'/>" +
                        "</div>";
                }
            });
            inputGroupElement.classList.add("suggestion-hover");
            suggestionElement.style.display = 'flex';
        }
    } else {
        if (type === "mobile") {
            mobileSugItemsElement.innerHTML = "";
        } else {
            suggestionElement.style.display = 'none';
            inputGroupElement.classList.remove("suggestion-hover")
        }
    }
}

function clickDeleteHistoryKeyword(keyword, type) {
    event.stopPropagation();
    const sugs = handleDeleteHistoryKeywords(keyword, suggestions);
    suggestions = sugs;
    renderSuggestionsWithData(type, sugs);
}

function searchWithValue(operationType, screenType, searchValue) {
    setTimeout(() => {
        let value;
        if (operationType === "sugLink") {
            value = searchValue
        } else if (operationType === "inputConfirm") {
            if (screenType === "mobile") {
                value = mobileSugInputElement.value
            } else if (screenType === "pc") {
                value = inputElement.value
            }
        }
        if (value) {
            window.location.assign(window.origin + "/search?q=" + encodeURIComponent(value));
            const hideFunction = screenType === "mobile" ? hideMobileSuggestion : hideSuggestion;
            setTimeout(hideFunction, 500);
            setHistoryKeywordsWithSearch(value);
        }
    }, 100);
}

function hideSuggestion() {
    suggestionElement.style.display = 'none';
    inputGroupElement.classList.remove("suggestion-hover")
}

function hideMobileSuggestion() {
    mobileSugElement.style.display = "none";
    mobileSugInputElement.value = "";
    mobileSugCloseElement.style.display = 'none';
    mobileSugItemsElement.innerHTML = "";
};
mobileSugInputElement.addEventListener("keydown", function (event) {
    if (event.keyCode == 13) {
        searchWithValue("inputConfirm", "mobile");
    }
});

function handleMobileSuggestionInput(params) {
    handleKeywordInput({type: "mobile"});
    if (mobileSugInputElement.value) {
        mobileSugCloseElement.style.display = "flex";
    } else {
        mobileSugCloseElement.style.display = 'none';
    }
}

mobileSugCloseElement.addEventListener("click", function () {
    mobileSugInputElement.value = "";
    mobileSugCloseElement.style.display = 'none';
    mobileSugInputElement.focus();
});
mobileSugBackElement.addEventListener("click", function () {
    hideMobileSuggestion();
});
moreToolsSwitchElement.addEventListener('click', function (event) {
    event.stopPropagation();
    if (moreToolsContainerElement.style.display === 'flex') {
        moreToolsContainerElement.style.display = 'none';
        moreToolsContainerElement.classList.remove('active');
    } else {
        moreToolsContainerElement.style.display = 'flex';
        toolsContainerElement.style.display = 'none';
        moreToolsContainerElement.classList.add('active')
    }
});
const _showKeywordSuggestion = localStorage.getItem("showKeywordSuggestion") !== "false";
const handleKeywordInput = async (options = {}) => {
    const {type} = options;
    const inputValue = type === "mobile" ? mobileSugInputElement.value : inputElement.value;
    if (_showKeywordSuggestion) {
        const sugs = await generateSuggestions(inputValue) || [];
        renderSuggestionsWithData(type, sugs);
        suggestions = sugs;
    }
}
inputElement.addEventListener('keydown', function (event) {
    if (event.keyCode == 13) {
        searchWithValue("inputConfirm", "pc")
    } else if (event.keyCode == 27) {
        suggestionElement.style.display = 'none';
        inputGroupElement.classList.remove("suggestion-hover");
    } else if (event.keyCode == 38) {
        if (suggestions.length !== 0) {
            if (_suggestionIndex === null || _suggestionIndex === 0) {
                _suggestionIndex = suggestions.length - 1;
            } else {
                _suggestionIndex = _suggestionIndex - 1;
            }
        }
    } else if (event.keyCode == 40) {
        if (suggestions.length !== 0) {
            if (_suggestionIndex === null || _suggestionIndex === suggestions.length - 1) {
                _suggestionIndex = 0;
            } else {
                _suggestionIndex = _suggestionIndex + 1;
            }
        }
    }
    if (_suggestionIndex !== null && suggestions[_suggestionIndex] && (event.keyCode == 38 || event.keyCode == 40)) {
        suggestions.forEach(item => {
            if (item === suggestions[_suggestionIndex]) {
                document.getElementById(item.type + item.text).style.backgroundColor = "#eee";
                if (document.getElementById(item.type + item.text + "delete-logo")) {
                    document.getElementById(item.type + item.text + "delete-logo").style.visibility = "visible";
                }
            } else {
                document.getElementById(item.type + item.text).style.backgroundColor = "#FFF";
                if (document.getElementById(item.type + item.text + "delete-logo")) {
                    document.getElementById(item.type + item.text + "delete-logo").style.visibility = "";
                }
            }
        });
        inputElement.value = suggestions[_suggestionIndex].text;
    }
});
if (!isSmallDevice) {
    recordNumberElement.style.display = "inline";
    operationContainerElement.style.display = "flex";
    inputElement.style.display = 'inline-block';
    inputDivElement.style.display = 'none';
    bottomRightBarElement.style.position = "absolute";
    bottomRightBarElement.style.right = "50px";
    inputElement.focus();
    privacyEntryElement.innerText = "隐私条约";
    let isTapping = false;
    inputElement.addEventListener("compositionstart", function () {
        isTapping = true;
    });
    inputElement.addEventListener("compositionend", function () {
        isTapping = false;
        handleHomeInput();
    });
    inputElement.addEventListener("input", function (event) {
        if (!isTapping) {
            debounce(() => {
                handleHomeInput();
            }, 300);
        }
    });
    window.addEventListener("click", () => {
        hideSuggestion();
        moreToolsContainerElement.style.display = 'none';
        toolsContainerElement.style.display = "none";
        moreToolsContainerElement.classList.remove('active');
        themeContainerElement.style.display = "none";
    });
    inputGroupElement.addEventListener('click', function () {
        event.stopPropagation();
    });
    suggestionElement.addEventListener('click', function () {
        event.stopPropagation();
    });
} else {
    inputElement.style.display = 'none';
    operationWrapElement.style.top = "70px";
    inputDivElement.style.display = 'inline-block';
    realTimeConElement.classList.add("mobile");
    inputGroupElement.addEventListener('click', function () {
        mobileSugElement.style.display = "block";
        mobileSugInputElement.focus();
    });
    textLabelElement.addEventListener('click', function () {
        mobileSugElement.style.display = "block";
        mobileSugInputElement.focus();
    });
    let isTapping = false;
    mobileSugInputElement.addEventListener("compositionstart", function () {
        isTapping = true;
    });
    mobileSugInputElement.addEventListener("compositionend", function () {
        isTapping = false;
        handleMobileSuggestionInput();
    });
    mobileSugInputElement.addEventListener("input", function () {
        if (!isTapping) {
            debounce(() => {
                handleMobileSuggestionInput();
            }, 300);
        }
    });
    mobileSugInputElement.addEventListener("focus", function () {
        handleMobileSuggestionInput();
    });
    window.addEventListener("click", () => {
        themeContainerElement.style.display = "none";
    });
}

function handleHomeInput() {
    handleKeywordInput();
    if (inputElement.value) {
        closeElement.style.display = 'inline-block';
    } else {
        closeElement.style.display = 'none';
    }
}

inputElement.addEventListener('click', function () {
    moreToolsContainerElement.style.display = 'none';
    moreToolsContainerElement.classList.remove('active');
    if (suggestionElement.style.display !== "flex") {
        handleKeywordInput();
    }
});
closeElement.addEventListener("click", function () {
    inputElement.value = '';
    closeElement.style.display = 'none';
    inputElement.focus();
    handleKeywordInput();
});
doggoToggleElement.addEventListener("click", function () {
    if (!doggoElement.checked) {
        doggoToggleElement.style.backgroundColor = "#1973E8";
        doggoToggleToolTipContainerElement.style.left = "66px";
        doggoToggleToolTipContainerElement.style.transition = "0.3s all ease-in-out";
    } else {
        doggoToggleElement.style.backgroundColor = "#CCC";
        doggoToggleToolTipContainerElement.style.left = "32px";
        doggoToggleToolTipContainerElement.style.transition = "0.3s all ease-in-out";
    }
    localStorage.setItem("enableBetaFunction", !doggoElement.checked);
});
let realTimeNews = [];
let realTimeOffset = 0;

async function clickRealTimeBtn() {
    if (realTimeConElement.style.display !== "block") {
        const _theme = localStorage.getItem("themeModeConfig");
        realTimeConElement.style.display = "block";
        setTimeout(() => {
            refreshImgElement.src = _theme === "dark" ? "https://cdn.fsofso.com/static/assets/theme/icon_refresh_dark.svg" : "https://cdn.fsofso.com/static/assets/icon_refresh.svg";
            refreshElement.style.display = "block";
            footerElement.style.display = "none";
        }, 300)
        let realTimeDisplay;
        if (realTimeNews.length === 0) {
            try {
                const result = await getRealTimeNews();
                const data = await result.json();
                realTimeNews = data.topics;
                realTimeDisplay = realTimeNews.slice(0, 7);
                realTimeOffset = 7;
            } catch (error) {
                isDev && console.log('error :>> ', error);
            }
        } else {
            if (realTimeNews.length - realTimeOffset < 7) {
                realTimeOffset = 0;
            }
            realTimeDisplay = realTimeNews.slice(realTimeOffset, realTimeOffset + 7);
            realTimeOffset += 7;
        }
        realTimeConElement.innerHTML = "";
        const imgSrc = _theme === "dark" ? "https://cdn.fsofso.com/static/assets/theme/home_real_time_news_dark.svg" : "https://cdn.fsofso.com/static/assets/home_real_time_news.svg";
        realTimeDisplay.forEach((item, index) => {
            realTimeConElement.innerHTML += "<div" + " onclick='clickRealTimeNewLink(" + ('"' + item.link + '"') + ")' class='real-time-item'>" +
                "<div class='real-time-content'>" +
                "<div class='real-time-icon'>" +
                "<img src=" + imgSrc + " />" +
                "</div>" +
                "<div class='real-time-tile'>" + item.title + "</div>" +
                "</div>" +
                "<div class='real-time-line'></div>" +
                "</div>";
        });
    } else {
        realTimeConElement.style.display = "none";
        refreshElement.style.display = "none";
        footerElement.style.display = "block";
    }
};

function clickRefreshBtn() {
    if (realTimeNews.length - realTimeOffset < 7) {
        realTimeOffset = 0;
    }
    const realTimeDisplay = realTimeNews.slice(realTimeOffset, realTimeOffset + 7);
    realTimeOffset += 7;
    realTimeConElement.innerHTML = "";
    const _theme = localStorage.getItem("themeModeConfig");
    const imgSrc = _theme === "dark" ? "https://cdn.fsofso.com/static/assets/theme/home_real_time_news_dark.svg" : "https://cdn.fsofso.com/static/assets/home_real_time_news.svg";
    realTimeDisplay.forEach((item, index) => {
        realTimeConElement.innerHTML += "<div" + " onclick='clickRealTimeNewLink(" + ('"' + item.link + '"') + ")' class='real-time-item'>" +
            "<div class='real-time-content'>" +
            "<div class='real-time-icon'>" +
            "<img src=" + imgSrc + " />" +
            "</div>" +
            "<div class='real-time-tile'>" + item.title + "</div>" +
            "</div>" +
            "<div class='real-time-line'></div>" +
            "</div>";
    });
};

function getRealTimeNews() {
    const HOST = "https://apis.hetaoapis.com/" // isDev ? "https://devapis.hetaoapis.com/" : "https://apis.hetaoapis.com/";
    const _engine = "search-engine-listing/v1/trending/topics" // (isDev ? "dev-" : "") + "search-engine-listing/v1/trending/topics";
    const FULL_HOST = HOST + _engine;
    return fetch(FULL_HOST, {
        method: "GET",
        headers: {
            "x-ht-env": "prod", // isDev ? "dev" : "prod",
            "x-ht-s": "search",
            "Access-Control-Allow-Origin": "*",
        },
    })
}

function clickRealTimeNewLink(link) {
    if (isSmallDevice) {
        window.location.assign(link);
    } else {
        window.open(link);
    }
}

const emailAccount = localStorage.getItem("emailAccount");
if (emailAccount) {
    userSmallAvatarElement.style.display = 'flex';
    userSmallAvatarElement.innerText = emailAccount.slice(0, 1).toUpperCase();
    userBigAvatarElement.innerText = emailAccount.slice(0, 1).toUpperCase();
    emailDetailElement.innerText = emailAccount;
} else {
    userSmallAvatarElement.style.display = 'none';
}
userSmallAvatarElement.addEventListener("click", function (event) {
    event.stopPropagation();
    const toolsElementDisplay = toolsContainerElement.style.display;
    if (toolsElementDisplay === 'none' || toolsElementDisplay === '') {
        toolsContainerElement.style.display = 'flex';
        moreToolsContainerElement.style.display = 'none';
    } else {
        toolsContainerElement.style.display = 'none';
    }
});

function initThemeMode() {
    const _theme = localStorage.getItem("themeModeConfig");
    if (_theme === "dark") {
        DarkReader.enable({
            brightness: 100,
            contrast: 90,
            sepia: 10
        });
    } else if (_theme === "reading") {
        document.body.style.background = "#FFFEFC";
        operationWrapElement.style.background = "#FFFEFC";
        realTimeConElement.style.background = "#FFFEFC";
    } else {
        localStorage.setItem("themeModeConfig", "light");
    }
    initThemeModeItemImg();
    initThemeModeItemStyle();
};

function clickNavToSetting() {
    // window.location.assign(window.origin + "/setting");
};

function initThemeModeItemImg() {
    const _theme = localStorage.getItem("themeModeConfig");
    if (_theme === "dark") {
        themeLightImgElement.src = "https://cdn.fsofso.com/static/assets/theme/sun_dark.svg";
        themeReadingImgElement.src = "https://cdn.fsofso.com/static/assets/theme/book_dark.svg";
        themeDarkImgElement.src = "https://cdn.fsofso.com/static/assets/theme/moon_dark.svg";
    } else {
        themeLightImgElement.src = "https://cdn.fsofso.com/static/assets/theme/sun_light.svg";
        themeReadingImgElement.src = "https://cdn.fsofso.com/static/assets/theme/book_light.svg";
        themeDarkImgElement.src = "https://cdn.fsofso.com/static/assets/theme/moon_light.svg";
    }
};

function initThemeModeItemStyle() {
    const _theme = localStorage.getItem("themeModeConfig");
    inputElement.style.backgroundColor = "#FFF";
    if (_theme === "light") {
        themeLightElement.classList.add("active");
    } else if (_theme === "reading") {
        themeReadingElement.classList.add("active");
        inputElement.style.backgroundColor = "#FFFEFC";
    } else if (_theme === "dark") {
        themeDarkElement.classList.add("active");
        themeDarkElement.classList.add("dark");
        themeReadingElement.classList.add("dark");
        themeLightElement.classList.add("dark");
        toggleTooltipTextElement.innerHTML = "实验室 <br />1. Quick Code Search 功能<br />2. 实验算法 Ranking (搜索结果排序)";
        inputGroupElement.classList.add("dark");
        suggestionElement.classList.add("dark");
    }
    ;
};

function clickThemeMode(theme) {
    const _theme = localStorage.getItem("themeModeConfig");
    if (_theme !== theme) {
        localStorage.setItem("themeModeConfig", theme);
        window.location.reload();
    }
};

function clickSettingBtn() {
    event.stopPropagation();
    if (themeContainerElement.style.display !== "block") {
        themeContainerElement.style.display = "block";
    } else {
        themeContainerElement.style.display = "none";
    }
};

function stopProp() {
    event.stopPropagation();
};

const _host = "https://" + "apis.hetaoapis.com" // (isDev ? "dev" : "") + "apis.hetaoapis.com";
fetch(_host + "/utils/v1/ip/city?plaintext=1", {mode: 'cors'}).then(async result => {
    const response = await result.json();
    if (response["country"]) {
        localStorage.setItem('ll', JSON.stringify(response));
    }
});