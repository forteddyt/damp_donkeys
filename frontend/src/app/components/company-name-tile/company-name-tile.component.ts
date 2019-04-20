import { Component, Input } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';

export interface editInterface {
    deleteCompanyHelper(companyName: string, new_jwt: string);
}

@Component({
  selector: 'app-company-name-tile',
  templateUrl: './company-name-tile.component.html',
  styleUrls: ['./company-name-tile.component.css']
})

export class CompanyNameTileComponent  {
  companyName
  stateData
  careerFairName
  deleteCompanyHelper

  public selfRef: CompanyNameTileComponent

  //interface for Parent-Child interaction
  public compInteraction: editInterface;

  constructor(private http: HttpClient, private router: Router) {  }

  genNewCode(event: any){
    this.genNewCodePromise().then(
      (val) => { // success
        this.stateData.jwt = val["jwt"] // update jwt
        alert("Company \"" + val["company_name"] + "\" has new login code: " + val["user_code"])
      },
      (err) => { // failure
        if (err.status == 401){ // invalid jwt for request
          alert("Session expired")
          this.router.navigateByUrl('/admin')
        } else {
          alert("Whoops, something broke... status code: " + err.status)
        }
      });
  }

  genNewCodePromise() {
    return this.http.put("https://csrcint.cs.vt.edu/api/reset_code?company_name=" + this.companyName + "&jwt=" + this.stateData.jwt, {observe: 'response'}).toPromise();
  }

  deleteCompany(event: any){
    this.deleteCompanyPromise().then(
      (val) => {
        console.log(val)
        this.compInteraction.deleteCompanyHelper(this.companyName, val.body["jwt"])
      },
      (err) => {
        if (err.status == 401){// invalid jwt for request
          alert("Session expired")
          this.router.navigateByUrl('/admin')
        } else if (err.status == 403){// Deleting this company will break db
          alert("Cannot delete companies with checked in interviews")
        } else {
          alert("Whoops, something broke... status code: " + err.status)
        }
      });
  }

  deleteCompanyPromise(){
    return this.http.delete("https://csrcint.cs.vt.edu/api/delete_company?company_name=" + this.companyName + "&career_fair_name=" + this.careerFairName + "&jwt=" + this.stateData.jwt, {observe: 'response'}).toPromise();
  }
}
