class GeneDrug < ApplicationRecord
	belongs_to :gene
	belongs_to :drug
	belongs_to :dataset
	belongs_to :tissue
end
