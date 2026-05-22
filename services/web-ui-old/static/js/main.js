htmx.onLoad(function(target) {
    // Ищем элементы внутри target (нового контента) или во всем документе
    const sidebar = document.querySelector('.sidebar');
    const toggleBtn = document.querySelector('#toggleSidebarBtn');
    const brand = document.querySelector('.sidebar-brand');
    const mainContent = document.querySelector('.main-content');

    if (!sidebar || !toggleBtn || !brand) return;

    // 1. Восстанавливаем состояние из localStorage
    const isCollapsed = localStorage.getItem('sidebarCollapsed') === 'true';
    if (isCollapsed) {
        sidebar.classList.add('collapsed');
        mainContent?.classList.add('expanded');
    } else {
        sidebar.classList.remove('collapsed');
        mainContent?.classList.remove('expanded');
    }

    // 2. Сворачивание/разворачивание (удаляем старый слушатель, если есть, и вешаем новый)
    toggleBtn.onclick = function(e) {
        e.stopPropagation();
        sidebar.classList.toggle('collapsed');
        mainContent?.classList.toggle('expanded');
        const collapsed = sidebar.classList.contains('collapsed');
        localStorage.setItem('sidebarCollapsed', collapsed);
    };

    // 3. Разворачивание по клику на логотип
    brand.onclick = function() {
        if (sidebar.classList.contains('collapsed')) {
            sidebar.classList.remove('collapsed');
            mainContent?.classList.remove('expanded');
            localStorage.setItem('sidebarCollapsed', 'false');
        }
    };
});

// document.addEventListener('DOMContentLoaded', function() {
//     const sidebar = document.querySelector('.sidebar');
//     const toggleBtn = document.getElementById('toggleSidebarBtn');
//     const brand = document.querySelector('.sidebar-brand');
//     const mainContent = document.querySelector('.main-content');
//
//     if (!toggleBtn || !brand) return;
//
//     // Восстанавливаем состояние из localStorage
//     const isCollapsed = localStorage.getItem('sidebarCollapsed') === 'true';
//     if (isCollapsed) {
//         sidebar.classList.add('collapsed');
//         mainContent.classList.add('expanded');
//     }
//
//     // Сворачивание/разворачивание по кнопке
//     toggleBtn.addEventListener('click', function(e) {
//         e.stopPropagation();
//         sidebar.classList.toggle('collapsed');
//         mainContent.classList.toggle('expanded');
//         const collapsed = sidebar.classList.contains('collapsed');
//         localStorage.setItem('sidebarCollapsed', collapsed);
//     });
//
//     // Разворачивание по клику на логотип (только если свёрнут)
//     brand.addEventListener('click', function() {
//         if (sidebar.classList.contains('collapsed')) {
//             sidebar.classList.remove('collapsed');
//             mainContent.classList.remove('expanded');
//             localStorage.setItem('sidebarCollapsed', 'false');
//         }
//     });
// });
//
// document.body.addEventListener('htmx:configRequest', (event) => {
//     const token = document.querySelector('meta[name="csrf-token"]').getAttribute('content');
//     event.detail.headers['X-CSRF-Token'] = token;
// });
