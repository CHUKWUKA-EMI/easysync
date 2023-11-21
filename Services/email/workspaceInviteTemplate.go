package emailservice

import "fmt"

// WorkspaceInviteTemplate for building workspace invitation email content
type WorkspaceInviteTemplate struct {
	WorkspaceName string `json:"workspaceName" binding:"required"`
	EmailSubject  string `json:"emailSubject" binding:"required"`
	InvitationURL string `json:"invitationUrl" binding:"required"`
}

// Build method dynamically generates the email HTML content
func (w *WorkspaceInviteTemplate) Build() string {
	return fmt.Sprintf(`
   <!DOCTYPE html>
   <html lang="en">
	 <head>
	   <meta charset="UTF-8" />
	   <meta http-equiv="X-UA-Compatible" content="IE=edge" />
	   <meta name="viewport" content="width=device-width, initial-scale=1.0" />
	   <title>%s</title>
	   <style>
		 .btn {
		   padding: 0.8rem;
		   border-radius: 0.5rem;
		   outline: none;
		   border: none;
		   font-weight: 500;
		   color: white !important;
		   background-color: #0b0b9f;
		   cursor: pointer;
		   font-size: 16px;
		   letter-spacing: 2px;
		   text-decoration: none;
		 }
		 .btn:hover {
		   background-color: #080871;
		 }
		 .link {
		   color: #0b0be3;
		   font-weight: bold;
		   text-decoration: none;
		 }
		 .link:hover {
		   text-decoration: underline;
		 }
	   </style>
	   <script>
		 function openUrl(url) {
		   window.open(url);
		 }
	   </script>
	 </head>
	 <body>
	   <div style="margin-right: auto; margin-left: auto">
	   <h1>Join %s on EasySync</h1>
		 <p>
		   <a href="%s" class="btn">
			 JOIN
		   </a>
		 </p>
         <footer>
		 <p>
		 EasySync is a collaboration and messaging tool for teams and organizations.
		 </p>
		 </footer>
	   </div>
	 </body>
   </html>   
   `, w.EmailSubject, w.WorkspaceName, w.InvitationURL)
}
