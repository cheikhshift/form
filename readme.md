# Go-Form
Create HTML forms with your structs and interfaces.

### Requirements
- [Go Server](http://golangserver.com)


### Install and import
Add this tag within your `.gxml` file.

	<import src="github.com/cheikhshift/form/gos.gxml" />	


### Configure
Set max upload size :	

	netform.MaxSize = 20 //mb
Set form token key :
	
	netform.FormKey = "a very very very very secret key"
Set input field class :
	
	netform.InputClass = ""
Set submit button class :
	
	netform.ButtonClass = ""
### How to use

Build a new form within a Golang server template:
							
	{{ Build $arg1 $arg2 $arg3 $arg4 .Session }}
			
Argument information :

- `$arg1` : Interface to build form with. Submit a variable with data to prepopulate form.
- `$arg2` :  Target URL to submit form to.
- `$arg3` : Method of form submission.
- `$arg4` : Call to action of form button.
- `.Session` :  Current user session. Must be passed to ensure secure communication.
		
Build within `<end>` tag :

	net_Build(param1 interface{}, param2 string, param3 string, param4 string, param5 *sessions.Session) string

Parameter information :

- `param1` : Interface to build form with. Submit a variable with data to prepopulate form.
- `param2` :  Target URL to submit form to.
- `param3` : Method of form submission.
- `param4` : Call to action of form button.
- `param5` :  Current user session. Must be passed to ensure secure communication.

### Server side validation
Please visit the GoValidator page for valuable tag information. [here](https://github.com/asaskevich/govalidator)

### How to process data.
Within your `<end>` tag use the following function to validate and convert the post body to the specified interface.

	var sampleform SampleForm
    err := netform.Form(r, &sampleform)

### Field types
List of field types with associated tag behavior.


#### 1. string
Display text input box.

Tag properties :
- title : title of field.
- placeholder : placeholder of field.

#### 2. int | float (any number)
Display number input box.

Tag properties :
- title : title of field.
- placeholder : placeholder of field.

#### 3. bool
Display Checkbox.

Tag properties :
- title : text blurb right of checkbox.

#### 4. File
Display file upload box. Use this field property with function `netform.Path`to get local filesystem path.

Tag properties :
- title : title of field.
- file : Mimetype of file to upload. 

#### 5. Paragraph
Display text area.

Tag properties :
- title : title of field.
- placeholder : placeholder of field.

#### 6. Date
Display date input.

Tag properties :
- title : title of field.
- placeholder : placeholder of field.

#### 7. Select
Display dropdown list.

Tag properties :
- title : title of field.
- placeholder : Prompt left of field.
- select : comma delimited choices of field.

#### 8. SelectMult
Display dropdown list with multiple selection support.

Tag properties :
- title : title of field.
- placeholder : Prompt left of field.
- select : comma delimited choices of field.

#### 9. Radio
Display radio input.

Tag properties :
- title : title of field.
- select : comma delimited choices of field.

#### 10. Email
Display email input.

Tag properties :
- title : title of field.
- placeholder : placeholder of field.


### Samples

Sample of GoS `<struct/>` with form tags set : 

	<struct name="SampleForm">
			TestField string `title:"Hi world!",valid:"unique",placeholder:"Testfield prompt"`
			Count int `placeholder:"Count"`
			Name string `valid:"required",title:"Input title"`
			FieldTwo netform.Radio `title:"Enter Email",valid:"email,unique,required",select:"blue,orange,red,green"`
			FieldF netform.Select `placeholder:"Prompt?",valid:"email,unique,required",select:"blue,orange,red,green"`
			Created netform.Date
			Text netform.Paragraph 	`title:"Enter a description."`
			Photo netform.File 	`file:"image/*"`
			Emal netform.Email
			Terms bool	`title:"Accept terms of use."`
	</struct>
