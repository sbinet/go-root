=======
go-root
=======

Experimental, pure-Go package to read ROOT files (and perhaps write
them out too), without having ROOT installed.

Installation
============

::

  $ go get bitbucket.org/binet/go-root/pkg/groot


Example
=======

An executable ``groot-ls`` is provided, which will recursively dump
the hierarchical content of a ``ROOT`` file:

::

  $ go get bitbucket.org/binet/go-root/cmd/groot-ls
  $ groot-ls -f my.esd.root
  :: groot-ls ::
  file: 'my.esd.root' (version=53005)
  / -> #11 key(s)
  key: name='##Shapes' title='##Shapes' type=TTree
  key: name='##Links' title='##Links' type=TTree
  key: name='##Params' title='##Params' type=TTree
  key: name='CollectionTree' title='CollectionTree' type=TTree
  key: name='CollectionTreeTPCnv::MuonMeasurements_tlp2' title='CollectionTreeTPCnv::MuonMeasurements_tlp2' type=TTree
  key: name='CollectionTreeInDet::Track_tlp1' title='CollectionTreeInDet::Track_tlp1' type=TTree
  key: name='CollectionTreeMuonCaloEnergyContainer_tlp1' title='CollectionTreeMuonCaloEnergyContainer_tlp1' type=TTree
  key: name='CollectionTreeAnalysis::JetTagInfo_tlp3' title='CollectionTreeAnalysis::JetTagInfo_tlp3' type=TTree
  key: name='POOLContainer' title='POOLContainer' type=TTree
  key: name='MetaData' title='MetaData' type=TTree
  key: name='MetaDataHdr' title='MetaDataHdr' type=TTree
  ::bye.

