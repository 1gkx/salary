// Добавляем метод в jQuery
(function ($) {
  $.fn.extend({
    message: function (type, text) {
      // Находим контейнер с сообщениями
      let messages = document.querySelector('.messages')
      // Если его нет, создаем 
      if (!messages) {
        messages = document.createElement('div')
        messages.classList.add('messages')
        document.body.appendChild(messages);
      }
      // Создаем сообщение
      let $msg = document.createElement('div'),
        header = type == true ? "Success" : "Error";
      $msg.className = "alert alert-dismissible fade show";
      type == true ? $msg.classList.add("alert-success") : $msg.classList.add("alert-danger")
      $msg.innerHTML = `
            <h4 class="alert-heading">${header}</h4>
            <p>${text}</p>
            <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
        `;
      messages.appendChild($msg);
      setTimeout(function () {
        $msg.remove()
      }, 3000)
    },
    fail: function (msg) {
      let alert = $(this).find('.alert')
      alert.text(msg)
      alert.removeClass('d-none')
      setTimeout(function () {
        alert.addClass('d-none')
      }, 5000)
    },
    auth(responce) {
      // form = this
      if (response.code > 400) {
        return $(this).fail(response.data.message)
      }
      // Смс верификация прошла успешно
      if (response.code == 200 && response.data.verify) return document.location.href = '/'

      if (response.code == 200 && response.data.auth) {
        this.addClass('d-none')
        $('.verify').removeClass('d-none')
      }
    }
  })
}(jQuery));

document.addEventListener("DOMContentLoaded", function () {

  // Форматируем номер телефона
  $('#phone').on('keypress', function (e) {
    IMask(this, {
      mask: '+{7}(000)000-00-00'
    });
  })
  $('#phone').on('blur', function (e) {
    IMask(this, {
      mask: '+{7}(000)000-00-00'
    });
  })

  // Переключатели типа базы данных
  $('input[type=radio]').on('change', function (e) {
    if (this.value == "sqlite3") {
      $('.database_server').addClass('d-none')
      $('.database_path').removeClass('d-none')
    } else {
      $('.database_path').addClass('d-none')
      $('.database_server').removeClass('d-none')
    }
  });

  // Сохранение данных пользователя
  $(".adduser").submit(function (event) {
    event.preventDefault();
    if (this.checkValidity() === false) return this.classList.add('was-validated');
    $.post(this.action, JSON.stringify($(this).serializeArray()))
      .done(responce => $().message(true, responce))
      .fail(error => $().message(false, error));
  });

  // Авторизация и верификация
  $(".login").submit(function (event) {
    let self = this;
    event.preventDefault();
    if (this.checkValidity() === false) return this.classList.add('was-validated');
    $.post(this.action, JSON.stringify($(this).serializeArray()))
      .done(responce => $(self).auth(responce))
      .fail(error => $(self).fail(error));
  });

  // Удаление пользователя
  $(".delete").on('click', function (e) {
    e.preventDefault();
    let $el = e.target.closest('.item');
    if (el && el.dataset.id) {
      $.post("/admin/users", JSON.stringify($el.dataset.id))
        .done(responce => $().message(true, responce))
        .fail(error => $().message(false, error));
    } else {
      $().message(false, "Не найдет id пользователя")
    }
  });

  // Сохранение настроект
  $(".settings").submit(function (event) {
    event.preventDefault();
    if (this.checkValidity() === false) return this.classList.add('was-validated');

    // TODO Придумать более красивый вариант
    let data = new FormData(this),
      map = {};
    data.forEach((value, key) => {
      let k = key.split("_")
      if (k.length > 1) {
        if (!map[k[0]]) map[k[0]] = {}
        map[k[0]][k[1]] = value
      } else {
        map[key] = value
      }
    })

    console.log(map)
    $.post(this.action, JSON.stringify(map))
      .done(responce => $().message(true, responce))
      .fail(error => $().message(false, error));
  });

  // // btn_approved.forEach((btn) => {
  //   btn.addEventListener('click', function (e) {
  //     let $el = e.target.closest('.item');
  //     let data = {
  //       "ClientID": $el.dataset.clientid,
  //       "DepositID": $el.dataset.depositid,
  //       "PhoneNumber": $el.querySelector(".PhoneNumber").value
  //     }
  //     fetch('/approve', {
  //       method: 'POST',
  //       body: JSON.stringify(data)
  //     }).then(async (r) => {
  //       let response = await r.json()
  //       showResponce(true, response.status)
  //     }).catch((err) => {
  //       showResponce(false, err)
  //     })
  //   });
  // });

  // Установка приложения
  var navListItems = $('div.setup-panel a'),
    allWells = $('.setup-content'),
    allNextBtn = $('.nextBtn'),
    allBackBtn = $('.backBtn');

  // Скрываем все блоки
  allWells.hide();

  // Обработчики нажатия переключателей шагов
  navListItems.click(function (e) {
    e.preventDefault();
    var $target = $($(this).attr('href')),
      $item = $(this);

    if (!$item.hasClass('disabled')) {
      navListItems.removeClass('btn-primary').addClass('btn-default');
      $item.addClass('btn-primary');
      allWells.hide();
      $target.show();
      $target.find('input:eq(0)').focus();
    }
  });

  // Обработчик кнопк Далее
  allNextBtn.click(function () {
    var curStep = $(this).closest(".setup-content"),
      curStepBtn = curStep.attr("id"),
      nextStepWizard = $('div.setup-panel a[href="#' + curStepBtn + '"]').next(),
      curInputs = curStep.find("input[type='text'],input[type='url']"),
      isValid = true;

    $(".form-control").removeClass("is-invalid");
    for (var i = 0; i < curInputs.length; i++) {
      if (!curInputs[i].validity.valid) {
        isValid = false;
        $(curInputs[i]).closest(".form-control").addClass("is-invalid");
      }
    }

    if (isValid)
      nextStepWizard.removeAttr('disabled').trigger('click');
  });

  // Обработчик кнопк Назад
  allBackBtn.click(function () {
    var curStep = $(this).closest(".setup-content"),
      curStepBtn = curStep.attr("id"),
      nextStepWizard = $('div.setup-panel a[href="#' + curStepBtn + '"]').prev(),
      curInputs = curStep.find("input[type='text'],input[type='url']");

    nextStepWizard.removeAttr('disabled').trigger('click');
  });

  // Показать стартовый блок
  $('div.setup-panel a.btn-primary').trigger('click');

  // Обработчик выбора драйвера БД
  $('.db-drivers').click(function (e) {
    let type = $(e.target).children('input').val();

    if (type == 'sqlite3') {
      $('.bd_server').hide()
      $('.bd_path').show()
    } else {
      $('.bd_server').show()
      $('.bd_path').hide()
    }
  });

  $('button.save').click(function(e) {
    e.preventDefault();
    console.log(e.target)
  });

});