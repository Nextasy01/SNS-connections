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
    result.from_sns = sns === 'toggleYouTube' ? 'YouTube' : 'Instagram'
    result.to_sns = result.from_sns === 'YouTube' ? 'Instagram' : 'YouTube'
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

function imporToDrive(sns, popoverEl){
    console.log(sns)
    fetch("/view/google/import", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(getTheVideos(sns))
    })
    popoverEl.hide();
}
