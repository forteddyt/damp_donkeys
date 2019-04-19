import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';


@Component({
  selector: 'app-company-name-tile',
  templateUrl: './company-name-tile.component.html',
  styleUrls: ['./company-name-tile.component.css']
})
export class CompanyNameTileComponent  {
  companyName
  stateData
  constructor(private http: HttpClient, private router: Router) { }

  genNewCode(event: any){
    console.log(this.stateData.jwt)
    this.genNewCodePromise().then(
      (val) => { // success
        this.stateData.jwt = val["jwt"] // update jwt
        alert("Company \"" + val["company_name"] + "\" has new login code: " + val["user_code"])
      },
      (err) => { // failure
        if (err.status == 401){ // invalid jwt for request
          this.router.navigateByUrl('/admin')
          alert("Session expired")
        } else {
          alert("Whoops, something broke... status code: " + err.status)
        }
      });
  }

  genNewCodePromise() {
    return this.http.put("https://csrcint.cs.vt.edu/api/reset_code?company_name=" + this.companyName + "&jwt=" + this.stateData.jwt, {observe: 'response'}).toPromise();
  }

  deleteCompany(event: any){
    alert("Deleted company (but not really yet)")
  }
}
