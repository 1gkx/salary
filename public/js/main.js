document.addEventListener("DOMContentLoaded", function () {

  let formSignIn = document.querySelector('.signin'),
    formVerify = document.querySelector('.verify'),
    buttons = document.querySelectorAll('.btn');

  buttons.forEach(function ($btn) {
    $btn.addEventListener("click", function (e) {
      e.preventDefault()

      let $form = e.target.form;
      url = "/" + e.target.dataset.url;
      // Валидация
      if ($form.checkValidity() === false) return $form.classList.add('was-validated');
      post(url, $form);
    })
  });
});


async function post(url, form) {
  fetch(url, {
    method: 'POST',
    body: new FormData(form)
  }).then(async function (r) {
    let response = await r.json()
  
    if (response.code > 400) {
      return showError(
        form.querySelector('.alert'),
        response.data.message
      )
    }
    // Смс верификация прошла успешно
    if (response.code == 200 && response.data.verify) return document.location.href = '/'

    if (response.code == 200 && response.data.auth) {
      document.querySelector('.signin').classList.add('hidden')
      document.querySelector('.verify').classList.remove('hidden')
      document.querySelector('.verify').classList.add('visible')
    }
  }).catch(function (err) {
    showError(
      form.querySelector('.alert'),
      err
    )
    console.error("error!", err)
  })
}

function showError($alert, message) {
  $alert.classList.remove('hidden')
  $alert.classList.add('visible')
  $alert.innerText = message
  setTimeout(function () {
    $alert.classList.remove('visible')
    $alert.classList.add('hidden')
    $alert.innerText = ''
  }, 10000)
}


function createMessge(text) {
  if (!document.querySelector('.messages')) {
    let divMsg = document.createElement('div')
    divMsg.classList.add('messages')
    document.body.appendChild(divMsg);
  }
  // получаем контейнер
  const messages = document.querySelector('.messages');

  const message = document.createElement('div');
  message.className = 'alert alert-success alert-dismissible fade show';
  message.innerHTML = `
      <h4 class="alert-heading">Успех!</h4>
      <p>${text}</p>
      <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
  `;
  messages.appendChild(message);
}