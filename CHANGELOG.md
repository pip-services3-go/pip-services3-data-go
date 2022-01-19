# <img src="https://uploads-ssl.webflow.com/5ea5d3315186cf5ec60c3ee4/5edf1c94ce4c859f2b188094_logo.svg" alt="Pip.Services Logo" width="200"> <br/> Persistence components for Golang Changelog

## <a name="1.1.10"></a> 1.1.10 (2022-01-19) 
### Bug Fixes
- Fix GetListByFilter method in MemoryPersistence.

## <a name="1.1.9"></a> 1.1.9 (2021-10-22) 
### Features
- Updated dependencies for updating object fields by names in tags.

## <a name="1.1.8"></a> 1.1.8 (2021-10-21) 
### Bug Fixes
- Updated dependencies for fix integer values converting.

## <a name="1.1.7"></a> 1.1.7 (2021-08-06) 
### Bug Fixes
- Fix error in GenerateObjectId when Id field exists in nested struct

## <a name="1.1.5"></a> 1.1.5 (2021-07-07) 
### Bug Fixes
- Fix error in GetPageByFilter when paging skip greater than length of items slice

## <a name="1.1.4"></a> 1.1.4 (2021-05-06) 
### Bug Fixes
- Fix error in GetListByFilter when filter function is nil

## <a name="1.1.3"></a> 1.1.3 (2021-04-30) 
### Bug Fixes
- Fix deadlock in DeleteByFilter, DeleteById, Update and UpdatePartially methods

## <a name="1.1.0"></a> 1.1.0 (2021-04-04) 
### Features
* Move on Go 1.16

### Bug Fixes
* Fix CloneObject method in Utils

## <a name="1.0.7"></a> 1.0.7 (2020-12-11) 

### Features
* Update dependencies

## <a name="1.0.6"></a> 1.0.6 (2020-10-15) 

### Bug Fixes
* Fix visibility of GetIndexById method in IdentifiableMemoryPersistence


## <a name="1.0.5"></a> 1.0.5 (2020-07-12) 

### Features
* Moved some CRUD operations from IdentifiableMemoryPersistence to MemoryPersistence


## <a name="1.0.4"></a> 1.0.4 (2020-05-19) 

### Features
* Added GetCountByFilter method in IdentiifiableMemoryPersistence


## <a name="1.0.1-1.0.3"></a> 1.0.1-1.0.3 (2020-01-28) 

### Features
* Relocated general methods to utility module

### Bug Fixes
* Fix work with pionter type
* Fix deadlock in MemoryPersistence.Clear method
* Fix check paging param in GetPageByFilter method

## <a name="1.0.0"></a> 1.0.0 (2020-01-28) 

Initial public release

### Features
* **persistence** is a basic persistence that can work with any object types and provides only minimal set of operations
