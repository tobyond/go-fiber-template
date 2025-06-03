// Fresh App JavaScript

// Form handling with loading states
document.addEventListener('DOMContentLoaded', function() {
  // Add loading states to forms
  const forms = document.querySelectorAll('form')
  
  forms.forEach(form => {
    form.addEventListener('submit', function(e) {
      const submitBtn = form.querySelector('button[type="submit"]')
      if (submitBtn) {
        submitBtn.disabled = true
        submitBtn.classList.add('loading')
        
        // Add spinner to button
        const originalText = submitBtn.textContent
        submitBtn.innerHTML = `
          <span class="spinner mr-2"></span>
          ${originalText}
        `
        
        // Re-enable if form validation fails
        setTimeout(() => {
          if (!form.checkValidity()) {
            submitBtn.disabled = false
            submitBtn.classList.remove('loading')
            submitBtn.textContent = originalText
          }
        }, 100)
      }
    })
  })
  
  // Auto-hide alerts after 5 seconds
  const alerts = document.querySelectorAll('.alert')
  alerts.forEach(alert => {
    if (alert.classList.contains('alert-success')) {
      setTimeout(() => {
        alert.style.transition = 'opacity 0.5s'
        alert.style.opacity = '0'
        setTimeout(() => alert.remove(), 500)
      }, 5000)
    }
  })
  
  // Add smooth transitions to interactive elements
  const interactiveElements = document.querySelectorAll('button, .btn, input, .form-input')
  interactiveElements.forEach(el => {
    el.style.transition = 'all 0.2s ease'
  })
})

// Export for potential use in other modules
export default {
  // Add any app-wide functionality here
  init() {
    console.log('Fresh app initialized!')
  }
}
