const toggleContainers = document.querySelectorAll('[id^=toggle]');
const itemCheckboxes = document.querySelectorAll('[id^=checkbox]');

var popoverTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="popover"]'))
console.log(popoverTriggerList)
var popoverList = popoverTriggerList.map(function (popoverTriggerEl) {
  return new bootstrap.Popover(popoverTriggerEl)
})

var closePop = true;
document.body.addEventListener('mousedown', function(event) {
  if (event.target.classList.contains('popover-body') || event.target.classList.contains('importToYouTube') ||
  event.target.classList.contains('importToInsta')) {
    closePop = false;
    console.log(event.target)
  } else {
    closePop = true;
  }
  });

// Hide all item checkboxes by default
itemCheckboxes.forEach(checkbox => {
  checkbox.style.display = 'none';
});

// Add toggle event listener to each container toggle
toggleContainers.forEach(container => {
  container.addEventListener('change', () => {
    const containerItems = container.closest('.container').querySelectorAll('[id^=checkbox]')
    containerItems.forEach(item => {
      item.addEventListener('change', () => {
        const numChecked = container.closest('.container').querySelectorAll('[id^=checkbox]:checked').length
        item.labels[0].style.display = item.checked ? '' : 'none'
        selectButton.disabled = numChecked === 0

      })
    });

    if (container.checked){
      containerItems.forEach(item => {
        item.style.display = '';
      });
      container.closest('.container').querySelector('.select-button').style.display = '';
    }
    else{
      containerItems.forEach(item => {
        item.style.display = 'none'
      })
      container.closest('.container').querySelector('.select-button').style.display = 'none';
    }
  })
  const selectButton = document.createElement('button');
  selectButton.textContent = 'Import';
  selectButton.classList.add('btn', (container.id === 'toggleYouTube' ? 'btn-dark' : 'btn-light'), 'select-button');
  selectButton.type = 'button';
  selectButton.setAttribute("data-bs-container", "body");
  selectButton.setAttribute("data-bs-toggle", "popover");
  selectButton.setAttribute("data-bs-placement", "right");
  selectButton.setAttribute("data-bs-trigger", "manual");
  container.parentNode.insertBefore(selectButton, container.nextSibling.nextSibling.nextSibling.nextSibling);
  selectButton.style.display = 'none';
  selectButton.disabled = true;
  
  var popover = new bootstrap.Popover(selectButton, {
    html: true,
    content: () => {
      const popoverContent = "<p style='margin-bottom: 0'>Where to import?</p>" +
      "<div class='row'>" +
      "<div class='col pt-2 YT' style='display: none'>" +
      "<button type='submit' form='import_form1' style='background-color: white;border: none;padding: 0 0 0 0;width: 40px;height: 32px;'>" + 
      "<img src='svg/youtube_icon.png' alt='Import to YouTube' class='importToYouTube zoom' style='display: inline'></button></div>" + 
      "<div class='col pt-2 Insta' style='display: none'>" +
      "<button type='submit' form='import_form2' style='background-color: white;border: none;padding: 0 0 0 0;width: 40px;height: 32px;'>" +
      "<img src='svg/instagram_icon.png' alt='Import to Instagram' class='importToInsta zoom' style='display: inline'></button></div>" +
      "</div>"

      return popoverContent;
    },
    sanitize: false
  })

  // selectButton.addEventListener('shown.bs.popover', () => {
  //   const popoverEl = bootstrap.Popover.getInstance(selectButton);
  //   Array.from(popoverEl.tip.childNodes[1].childNodes).forEach(btn => {
  //       if (btn.nodeName === 'P'){
  //         return;
  //       }
  //       btn.addEventListener('click', imporToDrive)
  //   });
  // })

  selectButton.addEventListener('click', (e) =>{
    const popoverEl = bootstrap.Popover.getInstance(selectButton);
    console.log("Import clicked");

    popoverEl.tip.id = container.id === 'toggleYouTube' ? 'popYouTube' : 'popInsta';
    console.log(popoverEl.tip.id);
    console.log(popoverEl.tip.childNodes[1])
    Array.from(popoverEl.tip.childNodes[1].childNodes).forEach(el => {
      console.log(el.childNodes)

      if (el.nodeName === 'P'){
        return;
      }
      console.log(popoverEl.tip.id)
      console.log(el.firstChild.firstChild.firstChild.className)
      if (popoverEl.tip.id === 'popYouTube'){
        console.log("In importYT popup, displaying Insta btn")
        el.lastChild.style.display = 'block'
        el.addEventListener('click', () => {
          imporToDrive('toggleYouTube')
        })
      }
      if (popoverEl.tip.id === 'popInsta' ){
        console.log("In importInsta popup, displaying YT btn")
        el.firstChild.style.display = 'block'
        el.addEventListener('click', () => {
          imporToDrive('toggleInsta')
        })
      }
    })
  });

  selectButton.addEventListener('blur', (e) => {
    if(closePop){
      const popoverEl = bootstrap.Popover.getInstance(selectButton);
      popoverEl.hide()
    }
    else{
      selectButton.focus()
    }
  });
  selectButton.addEventListener('focus', (e) => {
    if(closePop){
      const popoverEl = bootstrap.Popover.getInstance(selectButton);
      popoverEl.show()
    }
  });
});