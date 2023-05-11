function getTheLink(sns){
    let result = []
    const toggleContainer = document.getElementById(sns)
    const allChecked = toggleContainer.closest('.container').querySelectorAll('[id^=checkbox]:checked')
    allChecked.forEach(checked => {
        const videos = checked.closest('.col-6').querySelectorAll('.card-body');
        videos.forEach(video => {
            result.push(video.firstElementChild.children[0].src)
        })
    });
    return result
}

function getTheVideos(sns){
    let result = {}
    result.sns = sns
    result.url = []
    result.title = []
    result.videoId = []
    const re = /(?<=checkbox).*/gm
    const toggleContainer = document.getElementById(sns)
    const allChecked = toggleContainer.closest('.container').querySelectorAll('[id^=checkbox]:checked')
    allChecked.forEach(checked => {
        const videoid = checked.id.match(re)
        result.videoId.push(videoid[0])
        var videos = checked.closest('.col-6').querySelectorAll('h4');
        videos.forEach(video => {
            result.title.push(video.innerText)
        })
        if (sns == "toggleYouTube"){
            return
        }
        videos = checked.closest('.col-6').querySelectorAll('.card-body');
        videos.forEach(video => {
            result.url.push(video.firstElementChild.children[0].src)
        })

    });
    console.log(result)
    return result


}

function getTheTitle(sns){
    let result = []
    const re = /(?<=checkbox).*/gm
    const toggleContainer = document.getElementById(sns)
    const allChecked = toggleContainer.closest('.container').querySelectorAll('[id^=checkbox]:checked')
    allChecked.forEach(checked => {
        const videoid = checked.id.match(re)
        console.log(videoid)
        const videos = checked.closest('.col-6').querySelectorAll('h4');
        videos.forEach(video => {
            result.push(video.innerText)
        })
    });
    console.log(result)
    return result
}

function imporToDrive(sns){
    const xhr = new XMLHttpRequest();
    xhr.open("POST", "/view/google/import");
    xhr.setRequestHeader("Content-Type", "application/json; charset=UTF-8")
    console.log(sns)
    // const links = getTheLink(sns);
    // const titles = getTheTitle(sns);
    // const jsonBody = {};
    // jsonBody.url = links
    // jsonBody.title = titles
    // console.log(jsonBody)
    // fetch("/google/import", {
    //     method: "POST",
    //     headers: {
    //         "Content-Type": "application/json",
    //         "Authorization": `Bearer ${token}`
    //     },
    //     body: JSON.stringify(getTheVideos(sns))
    // })
    xhr.send(JSON.stringify(getTheVideos(sns)));
}
