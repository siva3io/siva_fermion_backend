package elastic

import (
	"context"
	"os"

	el "github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License v3.0 as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License v3.0 for more details.
You should have received a copy of the GNU Lesser General Public License v3.0
along with this program.  If not, see <https://www.gnu.org/licenses/lgpl-3.0.html/>.
*/
var (
	ELASTIC_INDEX_NAME = os.Getenv("ELASTIC_INDEX_NAME")
)

var Client *el.Client

func Init() {

	var err error
	Client, err = el.NewClient(
		el.SetURL(os.Getenv("ELASTIC_URL_1"), os.Getenv("ELASTIC_URL_2")),
		el.SetSniff(false),
		el.SetBasicAuth(os.Getenv("ELASTIC_USERNAME"), os.Getenv("ELASTIC_PASSWORD")),
	)
	if err != nil {
		// panic(err)
	}
}

func Insert(ctx context.Context, index string, log interface{}) error {

	if _, err := Client.Index().Index(ELASTIC_INDEX_NAME).
		Type("_doc").
		BodyJson(log).
		Do(ctx); err != nil {

		logrus.WithFields(logrus.Fields{
			"ElasticSearch": "cannot insert data",
			"Index":         index,
			"Data":          log,
		}).Error(err.Error())
		return err
	}

	return nil
}

func Get(ctx context.Context) (interface{}, error) {

	return nil, nil
}

func Update(ctx context.Context, index, ID string, update map[string]interface{}) error {

	if _, err := Client.Update().
		Index(index).
		Type("_doc").
		Id(ID).Doc(update).Do(ctx); err != nil {

		logrus.WithFields(logrus.Fields{
			"ElasticSearch": "cannot update data",
			"ID":            ID,
			"Index":         index,
			"Data":          update,
		}).Error(err.Error())
		return err
	}

	return nil
}

func Search(ctx context.Context, index string, searchSource *el.SearchSource) (*el.SearchResult, error) {

	results, err := Client.Search().
		Index(index).
		SearchSource(searchSource).
		Do(ctx)

	if err != nil {

		logrus.WithFields(logrus.Fields{
			"ElasticSearch": "cannot search data",
		}).Error(err.Error())

		return nil, err
	}

	return results, nil
}
