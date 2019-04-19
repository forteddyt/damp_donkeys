import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import * as jwt_decode from 'jwt-decode';
import { Router, ActivatedRoute, NavigationExtras } from '@angular/router';

@Component({
  selector: 'app-admin',
  templateUrl: './admin.component.html',
  styleUrls: ['./admin.component.css']
})
export class AdminComponent implements OnInit {
	code = ''
	constructor(private http: HttpClient, private router: Router, private route: ActivatedRoute) { }

	ngOnInit() {
		document.getElementById("admin_code").focus();
	}

	onKey(event: any)
	{
		this.code = event.target.value;
	}

	handleClick(event: any)
	{
		this.http.get("https://csrcint.cs.vt.edu/api/login?code=" + this.code, {observe: 'response'}).subscribe(
		    resp => {
		    	var jwt = resp.body["jwt"]
				var decoded = jwt_decode(jwt);

				// Only the user "admin" should be able to log in
				if(decoded["user"] == "admin"){
			    	this.router.navigateByUrl('/admin/nav');
				} else {
					alert("Invalid Code")					
				}

		    },
		    error => {
				alert("Invalid Code")
				console.log(error)
		    }
		  );
	}
}
