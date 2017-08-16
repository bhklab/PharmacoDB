# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# Note that this schema.rb definition is the authoritative source for your
# database schema. If you need to create the application database on another
# system, you should be using db:schema:load, not running all the migrations
# from scratch. The latter is a flawed and unsustainable approach (the more migrations
# you'll amass, the slower it'll run and the greater likelihood for issues).
#
# It's strongly recommended that you check this file into your version control system.

ActiveRecord::Schema.define(version: 20160924185806) do

  create_table "cell_tissues", force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=utf8" do |t|
    t.integer  "cell_id"
    t.integer  "tissue_id"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["cell_id"], name: "index_cell_tissues_on_cell_id", using: :btree
    t.index ["tissue_id"], name: "index_cell_tissues_on_tissue_id", using: :btree
  end

  create_table "cellosaurus", force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=utf8" do |t|
    t.text     "identifier",   limit: 65535
    t.text     "accession_id", limit: 65535
    t.text     "synonyms",     limit: 65535
    t.text     "creferences",  limit: 65535
    t.text     "ireferences",  limit: 65535
    t.text     "webpages",     limit: 65535
    t.text     "comments",     limit: 65535
    t.text     "strdata",      limit: 65535
    t.text     "diseases",     limit: 65535
    t.text     "sorigin",      limit: 65535
    t.text     "hierarchy",    limit: 65535
    t.text     "ofsi",         limit: 65535
    t.text     "sgofcell",     limit: 65535
    t.text     "category",     limit: 65535
    t.datetime "created_at",                 null: false
    t.datetime "updated_at",                 null: false
  end

  create_table "cells", primary_key: "cell_id", id: :integer, force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=utf8" do |t|
    t.string   "accession_id"
    t.text     "cell_name",    limit: 65535
    t.datetime "created_at",                 null: false
    t.datetime "updated_at",                 null: false
  end

  create_table "datasets", primary_key: "dataset_id", id: :integer, force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=utf8" do |t|
    t.text     "dataset_name", limit: 65535
    t.datetime "created_at",                 null: false
    t.datetime "updated_at",                 null: false
  end

  create_table "dose_responses", force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=utf8" do |t|
    t.integer  "experiment_id"
    t.float    "dose",          limit: 24
    t.float    "response",      limit: 24
    t.datetime "created_at",               null: false
    t.datetime "updated_at",               null: false
    t.index ["experiment_id"], name: "index_dose_responses_on_experiment_id", using: :btree
  end

  create_table "drugs", primary_key: "drug_id", id: :integer, force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=utf8" do |t|
    t.text     "drug_name",  limit: 65535
    t.datetime "created_at",               null: false
    t.datetime "updated_at",               null: false
  end

  create_table "experiments", primary_key: "experiment_id", id: :integer, force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=utf8" do |t|
    t.integer  "cell_id"
    t.integer  "drug_id"
    t.integer  "dataset_id"
    t.datetime "created_at", null: false
    t.datetime "updated_at", null: false
    t.index ["cell_id"], name: "index_experiments_on_cell_id", using: :btree
    t.index ["dataset_id"], name: "index_experiments_on_dataset_id", using: :btree
    t.index ["drug_id"], name: "index_experiments_on_drug_id", using: :btree
  end

  create_table "source_cell_names", force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=utf8" do |t|
    t.integer  "cell_id"
    t.integer  "source_id"
    t.text     "cell_name",  limit: 65535
    t.datetime "created_at",               null: false
    t.datetime "updated_at",               null: false
    t.index ["cell_id"], name: "index_source_cell_names_on_cell_id", using: :btree
    t.index ["source_id"], name: "index_source_cell_names_on_source_id", using: :btree
  end

  create_table "source_drug_names", force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=utf8" do |t|
    t.integer  "drug_id"
    t.integer  "source_id"
    t.text     "drug_name",  limit: 65535
    t.datetime "created_at",               null: false
    t.datetime "updated_at",               null: false
    t.index ["drug_id"], name: "index_source_drug_names_on_drug_id", using: :btree
    t.index ["source_id"], name: "index_source_drug_names_on_source_id", using: :btree
  end

  create_table "source_tissue_names", force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=utf8" do |t|
    t.integer  "tissue_id"
    t.integer  "source_id"
    t.text     "tissue_name", limit: 65535
    t.datetime "created_at",                null: false
    t.datetime "updated_at",                null: false
    t.index ["source_id"], name: "index_source_tissue_names_on_source_id", using: :btree
    t.index ["tissue_id"], name: "index_source_tissue_names_on_tissue_id", using: :btree
  end

  create_table "sources", primary_key: "source_id", id: :integer, force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=utf8" do |t|
    t.integer  "dataset_id"
    t.text     "source_name", limit: 65535
    t.datetime "created_at",                null: false
    t.datetime "updated_at",                null: false
    t.index ["dataset_id"], name: "index_sources_on_dataset_id", using: :btree
  end

  create_table "tissues", primary_key: "tissue_id", id: :integer, force: :cascade, options: "ENGINE=InnoDB DEFAULT CHARSET=utf8" do |t|
    t.text     "tissue_name", limit: 65535
    t.datetime "created_at",                null: false
    t.datetime "updated_at",                null: false
  end

end
