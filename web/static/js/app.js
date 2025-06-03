// web/assets/js/app.js
document.addEventListener("DOMContentLoaded", function() {
  const forms = document.querySelectorAll("form");
  forms.forEach((form) => {
    form.addEventListener("submit", function(e) {
      const submitBtn = form.querySelector('button[type="submit"]');
      if (submitBtn) {
        submitBtn.disabled = true;
        submitBtn.classList.add("loading");
        const originalText = submitBtn.textContent;
        submitBtn.innerHTML = `
          <span class="spinner mr-2"></span>
          ${originalText}
        `;
        setTimeout(() => {
          if (!form.checkValidity()) {
            submitBtn.disabled = false;
            submitBtn.classList.remove("loading");
            submitBtn.textContent = originalText;
          }
        }, 100);
      }
    });
  });
  const alerts = document.querySelectorAll(".alert");
  alerts.forEach((alert) => {
    if (alert.classList.contains("alert-success")) {
      setTimeout(() => {
        alert.style.transition = "opacity 0.5s";
        alert.style.opacity = "0";
        setTimeout(() => alert.remove(), 500);
      }, 5e3);
    }
  });
  const interactiveElements = document.querySelectorAll("button, .btn, input, .form-input");
  interactiveElements.forEach((el) => {
    el.style.transition = "all 0.2s ease";
  });
});
var app_default = {
  // Add any app-wide functionality here
  init() {
    console.log("Fresh app initialized!");
  }
};
export {
  app_default as default
};
//# sourceMappingURL=app.js.map
