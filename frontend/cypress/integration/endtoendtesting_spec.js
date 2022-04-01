describe("Find the index page", function(){
    it("Finds the index page", function(){
        cy.visit('http://localhost:4200/')
    })
})

describe("Find the login form", function(){
    it("Login Form working", function(){
        cy.contains('Login').click()
    })
})

describe("Fill the login form", function(){
 it('Interacting with text fields',()=>{
    cy.get('input[id="email"]')
      .type('test@gmail.com')
 .should('have.value','test@gmail.com')
 .get('input[id="password"]')
 .type('12345678')
 .should('have.value','12345678')
 .get('button[id="submit"]').click()
 .get('button[id="close"]').click()
 })
})

describe("Find the signup form", function(){
    it("Signup form working", function(){
        cy.contains('Sign up').click()
    })
})


describe("Fill the signup form", function(){
    it('Interacting with text fields',()=>{
       cy.get('input[id="firstname"]')
         .type('Mark')
    .should('have.value','Mark')
    .get('input[id="lastname"]')
    .type('Jacobs')
    .should('have.value','Jacobs')
    .get('input[id="useremail"]')
      .type('test@gmail.com')
 .should('have.value','test@gmail.com')
 .get('input[id="userpassword"]')
 .type('12345678')
 .should('have.value','12345678')
    .get('button[id="signupsubmit"]').click()
    .get('button[id="signupclose"]').click()
    })
   })