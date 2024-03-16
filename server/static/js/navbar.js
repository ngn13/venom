const links = document.querySelectorAll("#nav-link")

links.forEach(l => {
  const url = new URL(l.href)
  if (location.pathname.includes(url.pathname)) {
    l.style = "color: var(--white);"
  }
})
