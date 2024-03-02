/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/views/**/*.hbs"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
}