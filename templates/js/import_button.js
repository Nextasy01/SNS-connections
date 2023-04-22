document.addEventListener('DOMContentLoaded', function() {
    //const toggleCheckboxes = document.getElementById('toggleCheckboxes');
    const checkboxes = document.querySelectorAll('.checkboxArea');
    const switchGroup = document.querySelectorAll('.form-switch');
    const btnArea = document.querySelectorAll('.btnArea');
    
    
    const switchArea1 = document.getElementById('YouTubeImport');
    const switchArea2 = document.getElementById('InstaImport');


    switchGroup.forEach((switchGr) => {
      const switcher = switchGr.children[0];
      const ImportBtn = addButton(switcher.id);

      switcher.addEventListener('change', () => {
        btnArea.forEach((btn) => {
          if (switcher.id == 'InstaToggle'){
            switchAction(btn, switcher.checked, ImportBtn)
          }
          else{
            switchAction(btn, switcher.checked, ImportBtn)
          }
            
        })
      })
    });

    function addButton(platform) {
      const templateButton = document.createElement("button");

      templateButton.className = platform == 'InstaToggle' ? 'btn btn-light' : 'btn btn-dark';
      templateButton.textContent = 'Import';
      templateButton.type = 'button';
      templateButton.id = platform == 'InstaToggle' ? 'importInsta' : 'importYouTube';
      templateButton.disabled = true;
      
      return templateButton
    }

    function showOrHideButton(isChecked, area, btn){
      if(isChecked){
        area.appendChild(btn)
      }else{
        area.removeChild(btn)
      }
    }

    function switchAction(area, isChecked, btn){
      switch (area.id) {
        case 'InstaImport':
          showOrHideButton(isChecked, btn, ImportBtn);
          break;
        case 'YouTubeImport':
          showOrHideButton(isChecked, btn, ImportBtn);
          break;
        default:
          break;
      }
    }

    // toggleCheckboxes.addEventListener('change', () => {
    //     if(toggleCheckboxes.checked){
    //         switchArea1.appendChild(templateButton);
    //       }else{
    //         switchArea1.removeChild(templateButton);
    //       }

    //     checkboxes.forEach((checkboxArea) => {
    //       checkboxArea.style.display = toggleCheckboxes.checked ? 'flex' : 'none';

    //       const checkbox = checkboxArea.children[0].children[0];
    //       const label = checkboxArea.children[1].children[0];
          
    //       checkbox.addEventListener('change', () => {
    //         let isAnyCheckBoxChecked = false;

    //         checkboxes.forEach((checkboxArea2) => {
    //           if (checkboxArea2.children[0].children[0].checked){
    //             isAnyCheckBoxChecked = true;
    //           }
    //         });

    //         if (checkbox.checked && label.getAttribute('for') === checkbox.id){
    //           label.style.display = 'flex';
    //         }else{
    //           label.style.display = 'none';
    //         }

    //         if (isAnyCheckBoxChecked) {
    //          templateButton.disabled = false;
    //         }else{
    //           templateButton.disabled = true;
    //         }
    //       });
    //     });
    //   });
    });