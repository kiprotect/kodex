document.addEventListener('DOMContentLoaded', function () {

    var toggles = document.querySelectorAll('[data-toggle]').forEach(
        function(toggle){
            var target = toggle.dataset.target
            var el = document.getElementById(target)
            if (el === null)
                return
            var toggleClass = 'is-hidden'
            if (el.dataset !== undefined && el.dataset.toggleClass !== undefined)
                toggleClass = el.dataset.toggleClass
            if (el !== null){
                toggle.addEventListener('click', function(e){
                    el.classList.toggle(toggleClass)
                    toggle.classList.toggle('is-active')
                    e.preventDefault()
                })
            }    
        }
    )
});