import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { ActivatedRoute } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import * as jwt_decode from 'jwt-decode';

@Component({
  selector: 'app-company-unique',
  templateUrl: './company-unique.component.html',
  styleUrls: ['./company-unique.component.css']
})
export class CompanyUniqueComponent implements OnInit {
  code = "";
  validString: String = "abcdefghij";
  company_name = "Company_Name";
  json
  //json_middle

  constructor(private http: HttpClient, private route: ActivatedRoute, private router: Router) {
     //code = this.route.snapshot.paramMap.get('code');
  }

  ngOnInit() {
    this.code = this.route.snapshot.params.code;
    var student = this.http.get("https://csrcint.cs.vt.edu/api/login?password_hash=test").subscribe((res) => {
      console.log(res);
      this.json = res["jwt"];
      //this.json_middle = this.json.split('.')[1];
      var decoded = jwt_decode(this.json);
      this.company_name = decoded['user'];
    });
    //console.log(this.json_middle);
    //console.log(jwt_decode(this.json));

    //this.company_name = this.json.jwt;

    //this.code = this.route.snapshot.paramMap.get('code');
    //this.code = this.route.snapshot.params.code;
    //if (angular.equals(this.code, this.validString)) {}
    //else {
    //  this.router.navigate(['employers'], {});
    //}
  }

}
