class Tissue < ApplicationRecord
	has_many :source_tissue_names
	has_many :cell_tissues
	has_many :cells, through: :cell_tissues
	has_many :sources, through: :source_tissue_names
	has_many :experiments
end
