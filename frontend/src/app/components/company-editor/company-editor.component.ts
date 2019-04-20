import { Component, ComponentRef, OnInit, ComponentFactoryResolver, ViewChild, ViewContainerRef } from '@angular/core';
import { Router} from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { CompanyNameTileComponent } from '../company-name-tile/company-name-tile.component';
import { CareerFairTileComponent } from '../career-fair-tile/career-fair-tile.component'
import { StatisticsTileComponent } from '../statistics-tile/statistics-tile.component'

@Component({
  selector: 'app-company-editor',
  templateUrl: './company-editor.component.html',
  styleUrls: ['./company-editor.component.css']
})
export class CompanyEditorComponent implements OnInit {
	@ViewChild('companyInsert', { read: ViewContainerRef }) companyInsert: ViewContainerRef;
	@ViewChild('careerFairInsert', { read: ViewContainerRef }) careerFairInsert: ViewContainerRef;
	@ViewChild('statsInsert', { read: ViewContainerRef }) statsInsert: ViewContainerRef;

	stateData
	
	companyList = []
	companyRefs = []
	
	careerFairList = []
	careerFairRefs = []
	curCareerFair
	selectedCareerFair
	
	newCompanyName

	constructor(private http: HttpClient, private router: Router, private vcr: ViewContainerRef, private cfr: ComponentFactoryResolver) {
		this.stateData = this.router.getCurrentNavigation().extras.state;
		if(this.stateData == null || this.stateData.jwt == null){
			//console.log("Invalid state has been passed (or not passed at all);"+this.stateData);
			this.router.navigateByUrl('/admin');
		}
	}

	async ngOnInit(){
		let resp = await this.getCompanies();
		//console.log(resp.status)
		this.companyList = resp.body["company_list"];
		
		resp = await this.getCareerFairs();
		this.careerFairList = resp.body["career_fair_list"];
		this.curCareerFair = this.careerFairList[0];
		this.selectedCareerFair = this.curCareerFair
		this.stateData.jwt = resp.body["jwt"]; // update jwt


		this.loadCompanyComponents()
		this.loadCareerFairComponents()
		this.loadStatsComponents()
	}

	getCompanies() {
		return this.http.get("https://csrcint.cs.vt.edu/api/company_list", {observe: 'response'}).toPromise();
	}

	getCareerFairs() {
		return this.http.get("https://csrcint.cs.vt.edu/api/career_fair_list?jwt=" + this.stateData.jwt, {observe: 'response'}).toPromise();
	}

	clearStatsHelper(){
		this.statsInsert.clear();
	}

	addStatsHelper(){
		let cFactory = this.cfr.resolveComponentFactory(StatisticsTileComponent);
		let statsRef: ComponentRef<StatisticsTileComponent> = this.statsInsert.createComponent(cFactory);
		let statsComponent = statsRef.instance;

		statsComponent.stateData = this.stateData;
		statsComponent.careerFairName = this.selectedCareerFair;
	}

	loadStatsComponents(){
		this.clearStatsHelper();
		this.addStatsHelper();
	}

	clearCareerFairHelper(){
		this.careerFairInsert.clear();
		this.careerFairRefs = [];
	}

	selectCareerFairHelper(careerFairName, newCompanyList, newJWT){
		this.clearCareerFairHelper();
		this.clearCompanyHelper();

		this.selectedCareerFair = careerFairName;
		this.companyList = newCompanyList;
		this.stateData.jwt = newJWT;

		this.loadCompanyComponents()
		this.loadCareerFairComponents()
		this.loadStatsComponents()
	}

	addCareerFairHelper(careerFairName){
		let cFactory = this.cfr.resolveComponentFactory(CareerFairTileComponent);
		let careerFairRef: ComponentRef<CareerFairTileComponent> = this.careerFairInsert.createComponent(cFactory);
		let careerFairComponent = careerFairRef.instance;

		careerFairComponent.careerFairName = careerFairName;
		careerFairComponent.stateData = this.stateData;
		careerFairComponent.selfRef = careerFairComponent;
		
		// providing parent Component reference to get access to parent class methods
		careerFairComponent.compInteraction = this;

		// add reference for newly created component
		this.careerFairRefs.push(careerFairRef);
	}

	loadCareerFairComponents(){
		this.clearCareerFairHelper();
		for(var i in this.careerFairList){
			this.addCareerFairHelper(this.careerFairList[i]);
		}
	}

	clearCompanyHelper(){
		this.companyInsert.clear();
		this.companyRefs = [];
	}

	addCompanyHelper(companyName) {
		let cFactory = this.cfr.resolveComponentFactory(CompanyNameTileComponent);
		let companyRef: ComponentRef<CompanyNameTileComponent> = this.companyInsert.createComponent(cFactory);
		let companyComponent = companyRef.instance;

		companyComponent.companyName = companyName;
		companyComponent.stateData = this.stateData;
		companyComponent.selfRef = companyComponent;
		companyComponent.careerFairName = this.selectedCareerFair;
		
		// providing parent Component reference to get access to parent class methods
		companyComponent.compInteraction = this;

		// add reference for newly created component
		this.companyRefs.push(companyRef);
	}

	loadCompanyComponents() {
		this.clearCompanyHelper();
		for(var i in this.companyList){
			this.addCompanyHelper(this.companyList[i]);
		}
	}

	deleteCompanyHelper(companyName, newJWT) {
		if (this.companyInsert.length < 1){
			return
		}
		this.stateData.jwt = newJWT

		let componentRef = this.companyRefs.filter(x => x.instance.companyName == companyName)[0];
		let component: CompanyNameTileComponent = <CompanyNameTileComponent>componentRef.instance;

		let vcrIndex: number = this.companyInsert.indexOf(componentRef)

		// removing component from container
		this.companyInsert.remove(vcrIndex);

		// Filter out the now-delete company from both lists
		this.companyRefs = this.companyRefs.filter(x => x.instance.companyName !== companyName);
		this.companyList = this.companyList.filter(x => x !== companyName)
	}

	onKey(event: any)
	{
		this.newCompanyName = event.target.value;
	}

	addCompany(event: any){
		if (this.newCompanyName == null || this.newCompanyName == ""){
			alert("Invalid company name")
		} else if (this.companyList.includes(this.newCompanyName)){
			alert("Company \"" + this.newCompanyName + "\" is already registered for this career fair")
		} else {
			this.addCompanyPromise().then(
				(val) => { // success
					this.companyList.push(this.newCompanyName)
					
					// resort, ignoring case
					this.companyList.sort(function (a, b){
						return a.toLowerCase().localeCompare(b.toLowerCase());
					})
					
					this.stateData.jwt = val["jwt"] // update jwt
					
					alert("Company \"" + val["company_name"] + "\" has login code: " + val["user_code"])
					this.loadCompanyComponents()
				},
				(err) => { // failure
					if (err.status == 401){ // invalid jwt for request
						this.router.navigateByUrl('/admin')
					} else {
						alert("Whoops, something broke... status code: " + err.status)
					}
				});
		}
	}


	addCompanyPromise() {
		return this.http.put("https://csrcint.cs.vt.edu/api/add_company?company_name=" + this.newCompanyName + "&jwt=" + this.stateData.jwt, {observe: 'response'}).toPromise();
	}

	backRedirect(event: any){
		this.router.navigateByUrl('/admin/nav', { state: this.stateData })
	}
}
